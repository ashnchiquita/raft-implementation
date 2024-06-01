package core

import (
	"time"

	"tubes.sister/raft/node/application"
	"tubes.sister/raft/node/data"

	gRPC "tubes.sister/raft/gRPC/node/core"
)

type RaftNode struct {
	// Self data
	Address data.Address
	Type    NodeType

	// State
	Persistence data.Persistence
	Volatile    data.Volatile

	// App
	Application application.Application

	// Timer
	Timeout time.Duration

	// RPCs
	gRPC.UnimplementedAppendEntriesServiceServer
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

func (rn *RaftNode) SetTimeout() {
	if rn.Type == LEADER {
		rn.Timeout = HEARTBEAT_INTERVAL
	} else {
		rn.Timeout = RandomizeElectionTimeout()
	}
}
