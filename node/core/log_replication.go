package core

import (
	"context"
	"fmt"
	"log"
	"strings"

	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/data"
)

type resultStatus int

const (
	STAT_SUCCESS resultStatus = iota + 1
	STAT_FAILED
	STAT_TIMEOUT
	STAT_HEARTBEAT
)

type replicationResult struct {
	node   *data.ClusterData
	status resultStatus
}

func matchedInMajority(clusterList []data.ClusterData, population []data.Address, commitIndex int, excludedAddress data.Address) bool {
	majorityCount := 0
	excludedAddressInCluster := 0

	addressMap := make(map[data.Address]bool)

	for _, address := range population {
		addressMap[address] = true
	}

	for _, node := range clusterList {
		if _, ok := addressMap[node.Address]; !ok {
			continue
		}

		if node.Address.Equals(&excludedAddress) {
			excludedAddressInCluster++
			continue
		}
		// log.Printf("Node %v match index: %v (with pointer: %p)", node.Address, node.MatchIndex, &node)
		if node.MatchIndex > commitIndex {
			majorityCount++
		}
	}

	majorityThreshold := len(population)/2 + 1 // ceil(populationLength / 2)
	return majorityCount+excludedAddressInCluster >= majorityThreshold
}

func (rn *RaftNode) applyLogToApplication(matchIndex int) {
	rn.Volatile.CommitIndex = matchIndex
	currLog := rn.Persistence.Log[matchIndex]
	splitted := strings.Split(currLog.Value, ",")

	switch currLog.Command {
	case "SET":
		rn.Application.Set(splitted[0], splitted[1])
	case "DEL":
		rn.Application.Del(splitted[0])
	case "APPEND":
		rn.Application.Append(splitted[0], splitted[1])
	}
}

func (rn *RaftNode) startReplicatingLogs() {
	c := make(chan replicationResult)

	// logger
	memlogs := Cyan + "Memchange >> " + Reset

	// Initialize replication to all nodes in the cluster
	for idx, clusterData := range rn.Volatile.ClusterList {
		if clusterData.Address.Equals(&rn.Address) {
			continue
		}
		go rn.replicate(&rn.Volatile.ClusterList[idx], c)
	}
	// log.Printf("Current address: %v", rn.Address)

	isCommitted := false

	for _, clusterData := range rn.Volatile.ClusterList {
		if clusterData.Address.Equals(&rn.Address) {
			continue
		}

		result := <-c
		// log.Printf("Replication result: %v", result)
		// log.Printf("Receive pointer (startReplicatingLogs): %p", result.node)

		switch result.status {
		case STAT_TIMEOUT:
		case STAT_HEARTBEAT:
			continue
		case STAT_FAILED:
			result.node.NextIndex--
		case STAT_SUCCESS:
			result.node.MatchIndex++
			result.node.NextIndex++
		}

		if rn.Volatile.IsJointConsensus {
			majorityInOld := matchedInMajority(rn.Volatile.ClusterList, rn.Volatile.OldConfig, rn.Volatile.CommitIndex, rn.Address)
			majorityInNew := matchedInMajority(rn.Volatile.ClusterList, rn.Volatile.NewConfig, rn.Volatile.CommitIndex, rn.Address)

			if majorityInOld && majorityInNew && !isCommitted {
				// rn.Volatile.CommitIndex++
				rn.applyLogToApplication(result.node.MatchIndex)
				isCommitted = true

				committedEntry := rn.Persistence.Log[rn.Volatile.CommitIndex]
				if committedEntry.Command == "OLDNEWCONF" {
					log.Println(rn.Volatile.NewConfig)
					rn.logf(memlogs + "Committing oldnewconf")
					b, err := data.MarshallConfiguration(rn.Volatile.NewConfig)
					if err != nil {
						rn.logf(memlogs+"Invalid new config format; err: %s", err.Error())
					}

					marshalledNew := string(b)
					newConfig := data.LogEntry{Term: rn.Persistence.CurrentTerm, Command: "CONF", Value: marshalledNew}
					rn.Persistence.Log = append(rn.Persistence.Log, newConfig)
					rn.Persistence.Serialize()
					err = rn.ApplyNewClusterList(marshalledNew)

					if err != nil {
						log.Fatal(err.Error())
					}
				}
				log.Printf(Purple+"Log Rep >> "+Reset+"Commit index now: %d", rn.Volatile.CommitIndex)
			}
		} else {
			if !isCommitted && matchedInMajority(rn.Volatile.ClusterList, data.ClusterListToAddressList(rn.Volatile.ClusterList), rn.Volatile.CommitIndex, rn.Address) {
				// rn.Volatile.CommitIndex++
				rn.applyLogToApplication(result.node.MatchIndex)
				isCommitted = true
				log.Printf(Purple+"Log Rep >> "+Reset+"Commit index now: %d", rn.Volatile.CommitIndex)
			}

			committedEntry := rn.Persistence.Log[rn.Volatile.CommitIndex]
			if committedEntry.Command == "CONF" {
				amIInClusterList := false
				for _, node := range rn.Volatile.ClusterList {
					if node.Address.Equals(&rn.Address) {
						amIInClusterList = true
						break
					}
				}

				if !amIInClusterList {
					data.DisconnectClusterList(rn.Volatile.ClusterList)
					log.Fatalf("Exiting program, current server is no longer a part of the cluster")
				}
			}
		}
	}
}

func (rn *RaftNode) replicate(node *data.ClusterData, c chan replicationResult) {
	var (
		sendingEntries    []*gRPC.AppendEntriesArgs_LogEntry
		prevTerm          int
		resStatus         resultStatus
		isSendingHearbeat bool
	)

	if node.NextIndex > 0 {
		prevTerm = rn.Persistence.Log[node.NextIndex-1].Term
	} else {
		prevTerm = -1
	}
	// log.Printf("[CHECK] Address: %v; nextIndex: %d; loglen: %d", node.Address, node.NextIndex, len(rn.Persistence.Log))
	isSendingHearbeat = node.NextIndex >= len(rn.Persistence.Log)

	ctx, cancel := context.WithTimeout(context.Background(), HEARTBEAT_SEND_INTERVAL)
	defer cancel()

	if isSendingHearbeat {
		// Send heartbeat
		sendingEntries = make([]*gRPC.AppendEntriesArgs_LogEntry, 0)
	} else {
		sendingEntries = []*gRPC.AppendEntriesArgs_LogEntry{
			{
				Term:    int32(rn.Persistence.Log[node.NextIndex].Term),
				Command: rn.Persistence.Log[node.NextIndex].Command,
				Value:   rn.Persistence.Log[node.NextIndex].Value,
			},
		}
	}

	if !isSendingHearbeat {
		rn.logf(Purple+"AppendEntries >> "+Reset+"Sending Append entries to %v", node.Address)

		rn.announcef("====== REPRESENTATION ======")

		fmt.Println(Purple+"Term log		: ", int32(rn.Persistence.Log[node.NextIndex].Term))
		fmt.Println(Purple + "Value			: " + rn.Persistence.Log[node.NextIndex].Value)
		fmt.Println(Purple+"Curr term		: ", rn.Persistence.CurrentTerm)
		fmt.Println(Purple+"Commited index		:", rn.Volatile.CommitIndex)
		fmt.Println(Reset)
	} else {
		rn.logf(Green+"Heartbeat >> "+Reset+"Sending heartbeat to %v", node.Address)
	}

	reply, err := node.Client.AppendEntries(ctx, &gRPC.AppendEntriesArgs{
		LeaderAddress: &gRPC.AppendEntriesArgs_LeaderAddress{
			Ip:   rn.Address.IP,
			Port: int32(rn.Address.Port),
		},
		Term:         int32(rn.Persistence.CurrentTerm),
		Entries:      sendingEntries,
		PrevLogIndex: int32(node.NextIndex - 1),
		PrevLogTerm:  int32(prevTerm),
		LeaderCommit: int32(rn.Volatile.CommitIndex),
	})

	if err != nil {
		rn.logf(Purple+"AppendEntries >> "+Reset+"Error replicating logs to %v: %v", node.Address, err)
		c <- replicationResult{node: node, status: STAT_TIMEOUT}
		return
	} else if isSendingHearbeat {

		if !reply.Success {
			rn.logf("Heartbeat to %v did not succeed", node.Address)
			c <- replicationResult{node: node, status: STAT_FAILED}
			return
		} else {
			rn.logf("Successfully sent heartbeat to %v", node.Address)
		}
		c <- replicationResult{node: node, status: STAT_HEARTBEAT}
		return
	}

	if reply.Success {
		rn.logf("Successfully replicated logs to %v", node.Address)
		resStatus = STAT_SUCCESS
	} else {
		rn.logf("Failed to replicate logs to %v", node.Address)
		resStatus = STAT_FAILED
	}

	if reply.Term > int32(rn.Persistence.CurrentTerm) {
		rn.logf("Leader's term is outdated, reverting to follower")
		rn.setAsFollower()
	}

	// log.Printf("Receive pointer (replicate): %p", node)
	c <- replicationResult{node: node, status: resStatus}
}
