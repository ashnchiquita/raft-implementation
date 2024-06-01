package core

import (
	"fmt"
	"log"
	"net"
	"time"

	go_grpc "google.golang.org/grpc"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

// Initialize raft node to listen for GRPC requests from other nodes in the cluster
// and start the timer loop to monitor timeout
func (rn *RaftNode) InitializeServer() {
	go rn.startGRPCServer()
	rn.startTimerLoop()
}

// Starts the GRPC server
func (rn *RaftNode) startGRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", rn.Address.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []go_grpc.ServerOption
	grpcServer := go_grpc.NewServer(opts...)

	gRPC.RegisterHelloServer(grpcServer, rn)

	grpcServer.Serve(lis)
}

// Starts the timer loop based on the current node type
func (rn *RaftNode) startTimerLoop() {
	prev := time.Now()
	for {
		now := time.Now()
		elapsed := now.Sub(prev)
		rn.timeout = rn.timeout - elapsed
		prev = now

		// !: For testing only, remove these lines when timeout handling has been implemented
		log.Printf("Current timeout: %v", rn.timeout)
		time.Sleep(500 * time.Millisecond)

		if rn.timeout <= 0 {
			// TODO: Handle timeout based on the current node type
			break
		}
	}
}
