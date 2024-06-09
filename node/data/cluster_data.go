package data

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "tubes.sister/raft/gRPC/node/core"
)

type ClusterData struct {
	Address    Address
	NextIndex  int
	MatchIndex int
	Client     pb.AppendEntriesServiceClient
	conn       *grpc.ClientConn
}

func NewClusterData(address Address, nextIndex int, matchIndex int) *ClusterData {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", address.IP, address.Port), opts...)

	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	return &ClusterData{
		Address:    address,
		NextIndex:  nextIndex,
		MatchIndex: matchIndex,
		Client:     pb.NewAppendEntriesServiceClient(conn),
		conn:       conn,
	}
}

func ClusterListFromAddresses(addresses []Address, nextIndex int) []ClusterData {
	clusterDataList := make([]ClusterData, len(addresses))
	for idx, address := range addresses {
		clusterDataList[idx] = *NewClusterData(address, nextIndex, 0)
	}
	return clusterDataList
}

func DisconnectClusterList(clusterList []ClusterData) error {
	for _, clusterData := range clusterList {
		err := clusterData.conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
