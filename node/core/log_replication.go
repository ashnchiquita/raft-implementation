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

// func (rn *RaftNode) replicateLogs(start, end int) {
// 	// Ensures only leader can replicate its log entries
// 	if rn.Volatile.Type != data.LEADER {
// 		return
// 	}

// 	sendingEntries := make([]*gRPC.AppendEntriesArgs_LogEntry, end-start+1)
// 	for i := start; i <= end; i++ {
// 		entry := rn.Persistence.Log[i]
// 		sendingEntries = append(sendingEntries, &gRPC.AppendEntriesArgs_LogEntry{
// 			Term:    int32(entry.Term),
// 			Command: entry.Command,
// 			Value:   entry.Value,
// 		})
// 	}

// 	c := make(chan bool)
// 	for _, addr := range rn.Volatile.ClusterList {
// 		if addr.Address == rn.Address {
// 			continue
// 		}

// 		go func(addr data.Address, c chan bool) {
// 			var opts []grpc.DialOption
// 			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

// 			conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", addr.IP, addr.Port), opts...)
// 			if err != nil {
// 				log.Fatalf("Failed to dial server: %v", err)
// 			}
// 			defer conn.Close()

// 			client := gRPC.NewAppendEntriesServiceClient(conn)

// 			reply, err := client.AppendEntries(context.Background(), &gRPC.AppendEntriesArgs{
// 				LeaderAddress: &gRPC.AppendEntriesArgs_LeaderAddress{
// 					Ip:   rn.Address.IP,
// 					Port: int32(rn.Address.Port),
// 				},
// 				Term:         int32(rn.Persistence.CurrentTerm),
// 				Entries:      sendingEntries,
// 				PrevLogIndex: int32(start - 1),
// 				PrevLogTerm:  int32(rn.Volatile.LastApplied),
// 				LeaderCommit: int32(rn.Volatile.CommitIndex),
// 			})
// 			if err != nil {
// 				log.Printf("Error replicating logs to %v: %v", addr.Port, err)
// 				c <- false
// 			}

// 			c <- reply.Success
// 		}(addr.Address, c)
// 	}

// 	count := 0
// 	for i := 0; i < len(rn.Volatile.ClusterList)-1; i++ {
// 		if <-c {
// 			count++
// 		}
// 	}

//		log.Printf("Successfully replicated %d entries", count)
//	}
type replicationResult struct {
	node    *data.ClusterData
	success bool
}

func (rn *RaftNode) startReplicatingLogs() {
	c := make(chan replicationResult)

	// Initialize replication to all nodes in the cluster
	for _, node := range rn.Volatile.ClusterList {
		go rn.replicate(&node, c)
	}

	for {
		result := <-c
		log.Printf("Replication result: %v", result)

		if result.success {
			result.node.NextIndex++
			result.node.MatchIndex++
		}

		time.Sleep(500 * time.Millisecond)
		go rn.replicate(result.node, c)
	}
}

func (rn *RaftNode) replicate(node *data.ClusterData, c chan replicationResult) {
	var (
		sendingEntries []*gRPC.AppendEntriesArgs_LogEntry
		opts           []grpc.DialOption
		prevTerm       int
	)

	log.Printf("Replicating logs to %v", node.Address)

	if node.NextIndex >= len(rn.Persistence.Log) {
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

	if node.NextIndex > 0 {
		prevTerm = rn.Persistence.Log[node.NextIndex-1].Term
	} else {
		prevTerm = -1
	}

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", node.Address.IP, node.Address.Port), opts...)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	client := gRPC.NewAppendEntriesServiceClient(conn)

	reply, err := client.AppendEntries(context.Background(), &gRPC.AppendEntriesArgs{
		LeaderAddress: &gRPC.AppendEntriesArgs_LeaderAddress{
			Ip:   rn.Address.IP,
			Port: int32(rn.Address.Port),
		},
		Term:         int32(rn.Persistence.CurrentTerm),
		Entries:      sendingEntries,
		PrevLogIndex: int32(node.NextIndex - 1), //* NAON INI TEH ?????
		PrevLogTerm:  int32(prevTerm),
		LeaderCommit: int32(rn.Volatile.CommitIndex),
	})
	if err != nil {
		log.Printf("Error replicating logs to %v: %v", node.Address, err)
		c <- replicationResult{node: node, success: false}
	}

	c <- replicationResult{node: node, success: reply.Success}
}
