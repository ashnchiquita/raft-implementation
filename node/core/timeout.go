package core

import (
	"math/rand"
	"time"
)

const (
	HEARTBEAT_RECV_INTERVAL_MIN = 10 * time.Second // Previously: 100ms
	HEARTBEAT_RECV_INTERVAL_MAX = 15 * time.Second
	HEARTBEAT_SEND_INTERVAL     = 5 * time.Second // Previously: 100ms
	ELECTION_TIMEOUT_MIN        = 10 * time.Second
	ELECTION_TIMEOUT_MAX        = 15 * time.Second
	RPC_TIMEOUT                 = 120 * time.Millisecond
)

func RandomizeElectionTimeout() time.Duration {
	return time.Duration(int64(ELECTION_TIMEOUT_MIN) + rand.Int63n(int64(ELECTION_TIMEOUT_MAX-ELECTION_TIMEOUT_MIN)))
}

func RandomizeHeartbeatRecvInterval() time.Duration {
	return time.Duration(int64(HEARTBEAT_RECV_INTERVAL_MIN) + rand.Int63n(int64(HEARTBEAT_RECV_INTERVAL_MAX-HEARTBEAT_RECV_INTERVAL_MIN)))
}
