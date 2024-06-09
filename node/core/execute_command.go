package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/data"
)

func (rn *RaftNode) getLeaderClient() (gRPC.CmdExecutorClient, *grpc.ClientConn, error) {
	leaderAddr := rn.Volatile.LeaderAddress
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(leaderAddr.String(), opts...)
	client := gRPC.NewCmdExecutorClient(conn)
	return client, conn, err
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
			leaderClient, leaderConn, err := rn.getLeaderClient()
			defer leaderConn.Close()
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
	case "set": // write
		newLog := data.NewLogEntry(rn.Persistence.CurrentTerm, "SET", data.WithValue(fmt.Sprintf("%s,%s", msg.Vals[0], msg.Vals[1])))
		rn.Persistence.Log = append(rn.Persistence.Log, *newLog)
		rn.Persistence.Serialize()

		for rn.Volatile.CommitIndex < len(rn.Persistence.Log)-1 {
			continue
		}

		return &gRPC.ExecuteRes{Success: true, Value: "OK"}, nil
	case "get": // read
		value, ok := rn.Application.Get(msg.Vals[0])
		if !ok {
			return &gRPC.ExecuteRes{Success: false, Value: ""}, nil
		}
		return &gRPC.ExecuteRes{Success: true, Value: value}, nil
	case "strlen": // read
		length, ok := rn.Application.Strlen(msg.Vals[0])
		if !ok {
			return &gRPC.ExecuteRes{Success: false, Value: "0"}, nil
		}
		return &gRPC.ExecuteRes{Success: true, Value: strconv.Itoa(length)}, nil
	case "del": // write
		value, ok := rn.Application.Get(msg.Vals[0])
		if !ok {
			return &gRPC.ExecuteRes{Success: false, Value: "Key does not exist"}, nil
		}
		newLog := data.NewLogEntry(rn.Persistence.CurrentTerm, "DEL", data.WithValue(msg.Vals[0]))
		rn.Persistence.Log = append(rn.Persistence.Log, *newLog)
		rn.Persistence.Serialize()

		for rn.Volatile.CommitIndex < len(rn.Persistence.Log)-1 {
			continue
		}

		return &gRPC.ExecuteRes{Success: true, Value: value}, nil
	case "append": // write
		newLog := data.NewLogEntry(rn.Persistence.CurrentTerm, "APPEND", data.WithValue(fmt.Sprintf("%s,%s", msg.Vals[0], msg.Vals[1])))
		rn.Persistence.Log = append(rn.Persistence.Log, *newLog)
		rn.Persistence.Serialize()

		for rn.Volatile.CommitIndex < len(rn.Persistence.Log)-1 {
			continue
		}

		return &gRPC.ExecuteRes{Success: true, Value: "OK"}, nil
	case "getall": // read
		allData := rn.Application.GetAll()
		kvPairs := make([]map[string]string, 0, len(allData))
		for key, value := range allData {
			kvPairs = append(kvPairs, map[string]string{"key": key, "value": value})
		}
		kvPairsJson, err := json.Marshal(kvPairs)
		if err != nil {
			return nil, err
		}
		return &gRPC.ExecuteRes{Success: true, Value: string(kvPairsJson)}, nil
	case "delall": // write
		newLog := data.NewLogEntry(rn.Persistence.CurrentTerm, "DELALL")
		rn.Persistence.Log = append(rn.Persistence.Log, *newLog)
		rn.Persistence.Serialize()

		for rn.Volatile.CommitIndex < len(rn.Persistence.Log)-1 {
			continue
		}

		return &gRPC.ExecuteRes{Success: true, Value: "OK"}, nil
	case "config":
		rn.LeaderEnterJointConsensus(msg.Vals[0])
		for rn.Volatile.IsJointConsensus {
			continue
		}
		return &gRPC.ExecuteRes{Success: true, Value: "OK"}, nil
	case "ping":
		return &gRPC.ExecuteRes{Success: true, Value: "PONG"}, nil
	case "log":
		logStr, err := rn.Persistence.GetPrettyLog()
		if err != nil {
			return &gRPC.ExecuteRes{Success: false, Value: "Failed to get log"}, nil
		}
		return &gRPC.ExecuteRes{Success: true, Value: logStr}, nil
	default:
		return &gRPC.ExecuteRes{Success: false, Value: "Invalid command"}, nil
	}
}
