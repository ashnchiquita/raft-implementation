package data

import pb "tubes.sister/raft/gRPC/node/core"

type ClusterData struct {
	Address    Address
	NextIndex  int
	MatchIndex int
	Client     pb.AppendEntriesServiceClient
}
