package core

import (
	"math/rand"
	"time"
)

const (
	HEARTBEAT_RECV_INTERVAL = 60 * time.Second // Previously: 100ms
	HEARTBEAT_SEND_INTERVAL = 5 * time.Second  // Previously: 100ms
	ELECTION_TIMEOUT_MIN    = 30 * time.Second
	ELECTION_TIMEOUT_MAX    = 60 * time.Second
	RPC_TIMEOUT             = 50 * time.Millisecond
)

func RandomizeElectionTimeout() time.Duration {
	return time.Duration(int64(ELECTION_TIMEOUT_MIN) + rand.Int63n(int64(ELECTION_TIMEOUT_MAX-ELECTION_TIMEOUT_MIN)))
}
