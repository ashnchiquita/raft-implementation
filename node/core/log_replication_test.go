package core

import (
	"fmt"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/data"
)

func TestLogReplication(t *testing.T) {
	var opts []grpc.DialOption

	leaderAddr := data.NewAddress("localhost", 5000)

	leader := NewRaftNode(*leaderAddr)
	leader.Volatile.Type = data.LEADER

	leader.Persistence = *data.NewPersistence(*leaderAddr)
	newLog := data.NewLogEntry(leader.Persistence.CurrentTerm, "set", data.WithValue("1,2"))
	leader.Persistence.Log = append(leader.Persistence.Log, *newLog)

	followerAddr := data.NewAddress("localhost", 5001)

	follower := NewRaftNode(*followerAddr)
	follower.Persistence = *data.NewPersistence(*followerAddr)

	newServerResult := make(chan bool)
	go func(c chan bool) {
		var opts []grpc.ServerOption

		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", follower.Address.Port))
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer(opts...)
		gRPC.RegisterAppendEntriesServiceServer(grpcServer, follower)

		c <- true
		grpcServer.Serve(lis)
	}(newServerResult)

	if !<-newServerResult {
		t.Errorf("Failed to start new server")
	}

	targetNode := &data.ClusterData{Address: *data.NewAddress("localhost", 5001), MatchIndex: -1, NextIndex: 0}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", targetNode.Address.IP, targetNode.Address.Port), opts...)
	if err != nil {
		t.Errorf("Failed to dial server: %v", err)
	}
	targetNode.Client = gRPC.NewAppendEntriesServiceClient(conn)

	repResult := make(chan replicationResult)
	go leader.replicate(targetNode, repResult)
	<-repResult

	log.Printf("Follower log: %v", follower.Persistence.Log)
	if len(leader.Persistence.Log) != 1 {
		t.Errorf("Failed to add log to leader")
	}
	if len(follower.Persistence.Log) != 1 {
		t.Errorf("Follower log entries mismatch with leader")
	}
	if leader.Persistence.Log[0].Command != follower.Persistence.Log[0].Command {
		t.Errorf("Command mismatch")
	}
	if leader.Persistence.Log[0].Value != follower.Persistence.Log[0].Value {
		t.Errorf("Value mismatch")
	}
	if leader.Persistence.Log[0].Term != follower.Persistence.Log[0].Term {
		t.Errorf("Term mismatch")
	}
}
