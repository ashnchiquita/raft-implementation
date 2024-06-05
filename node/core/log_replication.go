package core

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	for {
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

		time.Sleep(500 * time.Millisecond)
		go rn.replicate(result.node, c)
	}
}

func (rn *RaftNode) replicate(node *data.ClusterData, c chan replicationResult) {
	var (
		sendingEntries    []*gRPC.AppendEntriesArgs_LogEntry
		opts              []grpc.DialOption
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

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", node.Address.IP, node.Address.Port), opts...)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	client := gRPC.NewAppendEntriesServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
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

	reply, err := client.AppendEntries(ctx, &gRPC.AppendEntriesArgs{
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
