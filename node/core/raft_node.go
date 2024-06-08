package core

import (
	"log"
	"time"

	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/application"
	"tubes.sister/raft/node/data"
)

type RaftNode struct {
	Address data.Address     // Address of the node (ip + port)
	timeout data.SafeTimeout // Current timeout value for the node

	// State
	Persistence data.Persistence
	Volatile    data.Volatile

	// App
	Application application.Application

	// RPCs
	gRPC.UnimplementedAppendEntriesServiceServer
	gRPC.UnimplementedHelloServer
	gRPC.UnimplementedRaftNodeServer
	gRPC.UnimplementedCmdExecutorServer

	// Channels
	electionInterrupt chan ElectionInterruptMsg
}

func NewRaftNode(address data.Address) *RaftNode {
	rn := &RaftNode{
		Address:           address,
		Application:       *application.NewApplication(),
		Volatile:          *data.NewVolatile(),
		electionInterrupt: make(chan ElectionInterruptMsg),
	}

	rn.Persistence = *data.InitPersistence(address)
	log.Printf("Logs: %v", rn.Persistence.Log)
	rn.resetTimeout()
	return rn
}

func (rn *RaftNode) resetTimeout() {
	rn.timeout.Mu.Lock()
	switch rn.Volatile.Type {
	case data.LEADER:
		rn.timeout.Value = HEARTBEAT_SEND_INTERVAL
	case data.CANDIDATE:
		rn.timeout.Value = RandomizeElectionTimeout()
	default:
		rn.timeout.Value = RandomizeHeartbeatRecvInterval()
	}
	rn.timeout.Mu.Unlock()
	log.Printf("Timeout reset to: %v", rn.timeout.Value)
}

func (rn *RaftNode) setTimoutSafe(val time.Duration) {
	rn.timeout.Mu.Lock()
	rn.timeout.Value = val
	rn.timeout.Mu.Unlock()
}
