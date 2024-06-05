package core

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/data"
)

func (rn *RaftNode) getLeaderClient() (gRPC.CmdExecutorClient, error) {
	leaderAddr := rn.Volatile.LeaderAddress
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(leaderAddr.String(), opts...)
	client := gRPC.NewCmdExecutorClient(conn)
	return client, err
}

func (rn *RaftNode) ExecuteCmd(ctx context.Context, msg *gRPC.ExecuteMsg) (*gRPC.ExecuteRes, error) {
	log.Printf("Received command: %v", msg.Cmd)

	if rn.Volatile.Type != data.LEADER {
		// Forward request to leader
		leaderAddr := rn.Volatile.LeaderAddress
		if leaderAddr.IsZeroAddress() {
			log.Printf("Try again later! No leader available.")
			return &gRPC.ExecuteRes{Success: false}, nil
		} else {
			leaderClient, err := rn.getLeaderClient()
			if err != nil {
				log.Printf("Error executing command: %v", err)
			} else {
				log.Printf("Successfully connected to leader")
				return leaderClient.ExecuteCmd(ctx, msg)
			}
		}
		return &gRPC.ExecuteRes{Success: false}, nil
	}

	switch msg.Cmd {
	case "set":
		newLog := data.NewLogEntry(rn.Persistence.CurrentTerm, "set", data.WithValue(fmt.Sprintf("%s,%s", msg.Vals[0], msg.Vals[1])))
		rn.Persistence.Log = append(rn.Persistence.Log, *newLog)
	}

	return &gRPC.ExecuteRes{Success: true}, nil
}
