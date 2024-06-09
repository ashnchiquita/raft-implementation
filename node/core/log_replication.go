package core

import (
	"context"
	"log"

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

func (rn *RaftNode) startReplicatingLogs() {
	c := make(chan replicationResult)

	// Initialize replication to all nodes in the cluster
	for idx, clusterData := range rn.Volatile.ClusterList {
		if clusterData.Address.Equals(&rn.Address) {
			continue
		}
		go rn.replicate(&rn.Volatile.ClusterList[idx], c)
	}
	log.Printf("Current address: %v", rn.Address)

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

			if majorityInOld && majorityInNew {
				rn.Volatile.CommitIndex++

				if entry := rn.Persistence.Log[rn.Volatile.CommitIndex]; entry.Command == "OLDNEWCONF" {
					log.Println(rn.Volatile.NewConfig)
					log.Printf("Memchange >> Committing oldnewconf")
					b, err := data.MarshallConfiguration(rn.Volatile.NewConfig)
					if err != nil {
						log.Printf("Memchange >> Invalid new config format; err: %s", err.Error())
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
				// return
			}
		} else {
			if matchedInMajority(rn.Volatile.ClusterList, data.ClusterListToAddressList(rn.Volatile.ClusterList), rn.Volatile.CommitIndex, rn.Address) {
				rn.Volatile.CommitIndex++
				// return
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
	log.Printf("Address: %v; nextIndex: %d; loglen: %d", node.Address, node.NextIndex, len(rn.Persistence.Log))
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

	log.Printf("Sending Append entries to %v", node.Address)
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
		log.Printf("Error replicating logs to %v: %v", node.Address, err)
		c <- replicationResult{node: node, status: STAT_TIMEOUT}
		return
	} else if isSendingHearbeat {
		if !reply.Success {
			log.Printf("Heartbeat to %v did not succeed", node.Address)
			c <- replicationResult{node: node, status: STAT_FAILED}
			return
		} else {
			log.Printf("Successfully sent heartbeat to %v", node.Address)
		}
		c <- replicationResult{node: node, status: STAT_HEARTBEAT}
		return
	}

	if reply.Success {
		log.Printf("Successfully replicated logs to %v", node.Address)
		resStatus = STAT_SUCCESS
	} else {
		log.Printf("Failed to replicate logs to %v", node.Address)
		resStatus = STAT_FAILED
	}

	if reply.Term > int32(rn.Persistence.CurrentTerm) {
		log.Printf("Leader's term is outdated, reverting to follower")
		rn.setAsCandidate()
	}

	// log.Printf("Receive pointer (replicate): %p", node)
	c <- replicationResult{node: node, status: resStatus}
}
