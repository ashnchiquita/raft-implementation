package data

import gRPC "tubes.sister/raft/gRPC/node/core"

type LeaderVolatile struct {
	AppendEntriesClients []gRPC.AppendEntriesServiceClient
	NextIndex            []int
	MatchIndex           []int
	// Volatile
}

// CONSTRUCTOR
func NewLeaderVolatile(commitIndex, lastApplied int) *LeaderVolatile {
	return &LeaderVolatile{
		AppendEntriesClients: []gRPC.AppendEntriesServiceClient{},
		NextIndex:            []int{},
		MatchIndex:           []int{},
		// Volatile: Volatile{
		// 	CommitIndex: commitIndex,
		// 	LastApplied: lastApplied,
		// },
	}
}
