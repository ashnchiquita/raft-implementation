package core

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/data"
)

// Initialize raft node to listen for GRPC requests from other nodes in the cluster
// and start the timer loop to monitor timeout
func (rn *RaftNode) InitializeServer() {
	go rn.startGRPCServer()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	for _, addr := range rn.Volatile.ClusterList {
		log.Printf("%s:%d", addr.IP, addr.Port)
		conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", addr.IP, addr.Port), opts...)

		if err != nil {
			log.Fatalf("Failed to dial server: %v", err)
		}

		client := pb.NewAppendEntriesServiceClient(conn)
		rn.LeaderVolatile.AppendEntriesClients = append(rn.LeaderVolatile.AppendEntriesClients, client)
		rn.LeaderVolatile.NextIndex = append(rn.LeaderVolatile.NextIndex, len(rn.Persistence.Log))
		rn.LeaderVolatile.MatchIndex = append(rn.LeaderVolatile.NextIndex, 0)
	}

	rn.startTimerLoop()
}

// Starts the GRPC server
func (rn *RaftNode) startGRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", rn.Address.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterHelloServer(grpcServer, rn)
	pb.RegisterAppendEntriesServiceServer(grpcServer, rn)

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
		// log.Printf("Current timeout: %v", rn.timeout)
		// time.Sleep(50 * time.Millisecond)

		if rn.timeout <= 0 {
			// TODO: Handle timeout based on the current node type

			switch rn.Volatile.Type {
			case data.LEADER:
				rn.heartbeat()
				rn.timeout = HEARTBEAT_SEND_INTERVAL
			case data.FOLLOWER:
				log.Println("start election")
				return
			}

			// break
		}
	}
}
