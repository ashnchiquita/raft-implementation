package core

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

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
		return &gRPC.ExecuteRes{Success: true, Value: "OK"}, nil
	case "get":
		var value string
		for _, logEntry := range rn.Persistence.Log {
			if logEntry.Command == "set" {
				keyAndValue := strings.Split(logEntry.Value, ",")
				if keyAndValue[0] == msg.Vals[0] {
					value = keyAndValue[1]
				}
			}
		}
		return &gRPC.ExecuteRes{Success: true, Value: value}, nil
	case "strlen":
		var value string
		for _, logEntry := range rn.Persistence.Log {
			if logEntry.Command == "set" {
				keyAndValue := strings.Split(logEntry.Value, ",")
				if keyAndValue[0] == msg.Vals[0] {
					value = keyAndValue[1]
				}
			}
		}
		length := len(value)
		return &gRPC.ExecuteRes{Success: true, Value: strconv.Itoa(length)}, nil
	case "del":
		var value string
		for i, logEntry := range rn.Persistence.Log {
			if logEntry.Command == "set" {
				keyAndValue := strings.Split(logEntry.Value, ",")
				if keyAndValue[0] == msg.Vals[0] {
					value = keyAndValue[1]
					rn.Persistence.Log[i].Command = "del"
				}
			}
		}
		return &gRPC.ExecuteRes{Success: true, Value: value}, nil
	case "append":
		var value string
		var keyExists bool
		for i, logEntry := range rn.Persistence.Log {
			if logEntry.Command == "set" {
				keyAndValue := strings.Split(logEntry.Value, ",")
				if keyAndValue[0] == msg.Vals[0] {
					value = keyAndValue[1] + msg.Vals[1]
					rn.Persistence.Log[i].Value = msg.Vals[0] + "," + value
					keyExists = true
				}
			}
		}
		if !keyExists {
			newLog := data.NewLogEntry(rn.Persistence.CurrentTerm, "set", data.WithValue(fmt.Sprintf("%s,%s", msg.Vals[0], msg.Vals[1])))
			rn.Persistence.Log = append(rn.Persistence.Log, *newLog)
		}
		return &gRPC.ExecuteRes{Success: true, Value: "OK"}, nil
	case "getall":
		var values []string
		for _, logEntry := range rn.Persistence.Log {
			if logEntry.Command == "set" {
				keyAndValue := strings.Split(logEntry.Value, ",")
				values = append(values, fmt.Sprintf("%s: %s", keyAndValue[0], keyAndValue[1]))
			}
		}
		return &gRPC.ExecuteRes{Success: true, Value: strings.Join(values, ", ")}, nil
	case "delall":
		for i, logEntry := range rn.Persistence.Log {
			if logEntry.Command == "set" {
				rn.Persistence.Log[i].Command = "del"
			}
		}
		return &gRPC.ExecuteRes{Success: true, Value: "OK"}, nil
	}

	return &gRPC.ExecuteRes{Success: true}, nil
}
