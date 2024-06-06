package core

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/data"
)

// Initialize raft node to listen for GRPC requests from other nodes in the cluster
// and start the timer loop to monitor timeout
func (rn *RaftNode) InitializeServer() {
	go rn.startGRPCServer()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// !: this code will make their own grpc client fail when called, which shouldn't happpen anyway
	for clusterIdx, clusterData := range rn.Volatile.ClusterList {
		// if clusterData.Address.Equals(&rn.Address) {
		// 	continue
		// }
		conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", clusterData.Address.IP, clusterData.Address.Port), opts...)

		if err != nil {
			log.Fatalf("Failed to dial server: %v", err)
		}

		client := gRPC.NewAppendEntriesServiceClient(conn)
		rn.Volatile.ClusterList[clusterIdx].Client = client
		rn.Volatile.ClusterList[clusterIdx].NextIndex = len(rn.Persistence.Log)
		rn.Volatile.ClusterList[clusterIdx].MatchIndex = 0
	}

	rn.startTimerLoop()
}

func (rn *RaftNode) InitializeAsLeader() {
	rn.Volatile.Type = data.LEADER
	rn.resetTimeout()
	go rn.startReplicatingLogs()
}

// Starts the GRPC server
func (rn *RaftNode) startGRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", rn.Address.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	gRPC.RegisterHelloServer(grpcServer, rn)
	gRPC.RegisterAppendEntriesServiceServer(grpcServer, rn)
	gRPC.RegisterCmdExecutorServer(grpcServer, rn)

	grpcServer.Serve(lis)
}

// Starts the timer loop based on the current node type
func (rn *RaftNode) startTimerLoop() {
	prev := time.Now()
	for {
		now := time.Now()
		elapsed := now.Sub(prev)
		rn.setTimoutSafe(rn.timeout.Value - elapsed)
		prev = now

		// !: For testing only, remove these lines when timeout handling has been implemented
		log.Printf("Current timeout: %v", rn.timeout.Value)
		time.Sleep(500 * time.Millisecond)

		if rn.timeout.Value <= 0 {
			switch rn.Volatile.Type {
			case data.LEADER:
				rn.startReplicatingLogs()
				rn.resetTimeout()
			case data.FOLLOWER:
				go rn.startElection()
			case data.CANDIDATE:
				rn.electionInterrupt <- ELECTION_TIMEOUT
			}
			log.Printf("Timeout occurred for node %v", rn.Address)
		}
	}
}
