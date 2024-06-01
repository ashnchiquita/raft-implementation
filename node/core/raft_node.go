package core

import (
	"fmt"
	"log"
	"net"
	"time"

	go_grpc "google.golang.org/grpc"
	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/application"
	"tubes.sister/raft/node/data"
)

type RaftNode struct {
	// Self data
	Address data.Address

	// State
	Persistence data.Persistence
	Volatile    data.Volatile

	// App
	Application application.Application

	// Timer
	Timeout time.Duration
	gRPC.UnimplementedHelloServer
}

func NewRaftNode(address data.Address, app application.Application) *RaftNode {
	rn := &RaftNode{
		Address:     address,
		Application: app,
		Volatile:    *data.NewVolatile(),
	}

	// TODO: try to load from file
	rn.Persistence = *data.NewPersistence()

	return rn
}

// Starts the node as GRPC server
func (rn *RaftNode) StartServer() {

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", rn.Address.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []go_grpc.ServerOption
	grpcServer := go_grpc.NewServer(opts...)

	gRPC.RegisterHelloServer(grpcServer, rn)

	grpcServer.Serve(lis)
}

func (rn *RaftNode) SetTimeout() {
	if rn.Volatile.Type == data.LEADER {
		rn.Timeout = HEARTBEAT_INTERVAL
	} else {
		rn.Timeout = RandomizeElectionTimeout()
	}
}
