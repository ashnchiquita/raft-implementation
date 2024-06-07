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

func (rn *RaftNode) startReplicatingLogs() {
	c := make(chan replicationResult)

	// Initialize replication to all nodes in the cluster
	for i := 1; i < len(rn.Volatile.ClusterList); i++ {
		// log.Printf("Send pointer: %p", &rn.Volatile.ClusterList[i])
		go rn.replicate(&rn.Volatile.ClusterList[i], c)
	}

	for range rn.Volatile.ClusterList[1:] {
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

			majorityCount := 0
			for i := 1; i < len(rn.Volatile.ClusterList); i++ {
				node := rn.Volatile.ClusterList[i]
				// log.Printf("Node %v match index: %v (with pointer: %p)", node.Address, node.MatchIndex, &node)
				if node.MatchIndex > rn.Volatile.CommitIndex {
					majorityCount++
				}
			}

			// log.Printf("Majority count: %v", majorityCount)
			if majorityCount+1 > len(rn.Volatile.ClusterList[1:])/2 {
				rn.Volatile.CommitIndex++
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
