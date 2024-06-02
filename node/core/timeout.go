package core

import (
	"math/rand"
	"time"
)

const (
	HEARTBEAT_RECV_INTERVAL = 60 * time.Second // Previously: 100ms
	HEARTBEAT_SEND_INTERVAL = 100 * time.Millisecond
	ELECTION_TIMEOUT_MIN    = 150 * time.Millisecond
	ELECTION_TIMEOUT_MAX    = 300 * time.Millisecond
)

func RandomizeElectionTimeout() time.Duration {
	return time.Duration(int64(ELECTION_TIMEOUT_MIN) + rand.Int63n(int64(ELECTION_TIMEOUT_MAX-ELECTION_TIMEOUT_MIN)))
}
