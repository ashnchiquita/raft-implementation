package core

import (
	"context"
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

func matchedInMajority(clusterList []data.ClusterData, commitIndex int, excludedAddress data.Address) bool {
	majorityCount := 0
	excludedAddressInCluster := 0
	for _, node := range clusterList {
		if node.Address.Equals(&excludedAddress) {
			excludedAddressInCluster++
			continue
		}
		// log.Printf("Node %v match index: %v (with pointer: %p)", node.Address, node.MatchIndex, &node)
		if node.MatchIndex > commitIndex {
			majorityCount++
		}
	}

	majorityThreshold := (len(clusterList) + 1) / 2 // ceil(clusterLength / 2)
	return majorityCount+excludedAddressInCluster >= majorityThreshold
}

func (rn *RaftNode) startReplicatingLogs() {
	c := make(chan replicationResult)

	// Initialize replication to all nodes in the cluster
	for _, clusterData := range rn.Volatile.ClusterList {
		// log.Printf("Send pointer: %p", &rn.Volatile.ClusterList[i])
		if clusterData.Address.Equals(&rn.Address) {
			continue
		}
		go rn.replicate(&clusterData, c)
	}

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
	}

	if rn.Volatile.IsJointConsensus {
		majorityInOld := matchedInMajority(rn.Volatile.OldClusterList, rn.Volatile.CommitIndex, rn.Address)
		majorityInNew := matchedInMajority(rn.Volatile.ClusterList, rn.Volatile.CommitIndex, rn.Address)

		if majorityInOld && majorityInNew {
			rn.Volatile.CommitIndex++
		}

		if entry := rn.Persistence.Log[rn.Volatile.CommitIndex]; entry.Command == "OLDNEWCONF" {
			marshalledNew := strings.Split(entry.Value, ",")[1]
			newConfig := data.LogEntry{Term: rn.Persistence.CurrentTerm, Command: "CONF", Value: marshalledNew}
			rn.Persistence.Log = append(rn.Persistence.Log, newConfig)
			err := rn.ApplyNewClusterList(marshalledNew)

			if err != nil {
				log.Fatal(err.Error())
			}
		}
	} else {
		if matchedInMajority(rn.Volatile.ClusterList, rn.Volatile.CommitIndex, rn.Address) {
			rn.Volatile.CommitIndex++
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
		log.Printf("Successfully sent heartbeat to %v", node.Address)
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
	// log.Printf("Receive pointer (replicate): %p", node)
	c <- replicationResult{node: node, status: resStatus}
}
