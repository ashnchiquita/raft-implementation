package core

import (
	"time"

	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/application"
	"tubes.sister/raft/node/data"
)

type RaftNode struct {
	Address data.Address  // Address of the node (ip + port)
	timeout time.Duration // Current timeout value for the node

	// State
	Persistence data.Persistence
	Volatile    data.Volatile

	// App
	Application application.Application

	// RPCs
	gRPC.UnimplementedAppendEntriesServiceServer
	gRPC.UnimplementedHelloServer
	gRPC.UnimplementedCmdExecutorServer
}

func NewRaftNode(address data.Address, app application.Application) *RaftNode {
	rn := &RaftNode{
		Address:     address,
		Application: app,
		Volatile:    *data.NewVolatile(),
	}

	// TODO: try to load from file
	rn.Persistence = *data.NewPersistence()

	rn.setTimeout()
	return rn
}

func (rn *RaftNode) setTimeout() {
	switch rn.Volatile.Type {
	case data.LEADER:
		rn.timeout = HEARTBEAT_SEND_INTERVAL
	case data.CANDIDATE:
		rn.timeout = RandomizeElectionTimeout()
	default:
		rn.timeout = HEARTBEAT_RECV_INTERVAL
	}
}
