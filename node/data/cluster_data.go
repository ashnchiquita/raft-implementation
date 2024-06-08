package data

import pb "tubes.sister/raft/gRPC/node/core"

type ClusterData struct {
	Address    Address
	NextIndex  int
	MatchIndex int
	Client     pb.AppendEntriesServiceClient
}

func NewClusterData(address Address, nextIndex int, matchIndex int) *ClusterData {
	return &ClusterData{
		Address:    address,
		NextIndex:  nextIndex,
		MatchIndex: matchIndex,
	}
}

func ClusterListFromAddresses(addresses []Address, nextIndex int) []ClusterData {
	clusterDataList := make([]ClusterData, len(addresses))
	for _, address := range addresses {
		clusterDataList = append(clusterDataList, *NewClusterData(address, nextIndex, 0))
	}
	return clusterDataList
}
