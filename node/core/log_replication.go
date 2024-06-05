package core

import (
	"context"
	"log"

	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/data"
)

type replicationResult struct {
	node    *data.ClusterData
	success string
}

func (rn *RaftNode) startReplicatingLogs() {
	c := make(chan replicationResult)

	// Initialize replication to all nodes in the cluster
	for _, node := range rn.Volatile.ClusterList[1:] {
		go rn.replicate(&node, c)
	}

	for i := 0; i < len(rn.Volatile.ClusterList[1:]); i++ {
		result := <-c
		log.Printf("Replication result: %v", result)

		switch result.success {
		case "timeout":
		case "hearbeat":
			continue
		case "failed":
			result.node.NextIndex--
		case "success":
			result.node.MatchIndex++
			result.node.NextIndex++
		}
	}
}

func (rn *RaftNode) replicate(node *data.ClusterData, c chan replicationResult) {
	var (
		sendingEntries    []*gRPC.AppendEntriesArgs_LogEntry
		prevTerm          int
		repResult         string
		isSendingHearbeat bool
	)

	log.Printf("Replicating logs to %v", node.Address)

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
		c <- replicationResult{node: node, success: "timeout"}
		return
	} else if isSendingHearbeat {
		c <- replicationResult{node: node, success: "heartbeat"}
		return
	}

	if reply.Success {
		repResult = "success"
	} else {
		repResult = "failed"
	}
	c <- replicationResult{node: node, success: repResult}
}
