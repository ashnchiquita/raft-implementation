package core

import (
	"context"
	"fmt"
	"log"

	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/data"
)

func (rn *RaftNode) ExecuteCmd(ctx context.Context, msg *gRPC.ExecuteMsg) (*gRPC.ExecuteRes, error) {
	log.Printf("Received command: %v", msg.Cmd)

	// TODO: Implement request forwarding to leader
	if rn.Volatile.Type != data.LEADER {
		return &gRPC.ExecuteRes{Success: false}, nil
	}

	switch msg.Cmd {
	case "set":
		newLog := data.NewLogEntry(rn.Persistence.CurrentTerm, "set", data.WithValue(fmt.Sprintf("%s,%s", msg.Vals[0], msg.Vals[1])))
		rn.Persistence.Log = append(rn.Persistence.Log, *newLog)
		log.Printf("New msg: %v", newLog)
	}

	return &gRPC.ExecuteRes{Success: true}, nil
}
