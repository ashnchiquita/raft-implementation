package application

import (
	"context"
	"encoding/json"
	"net/http"

	"tubes.sister/raft/client/http/utils"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

type DelAllResponse utils.KeyValResponse

func DelAll(client gRPC.CmdExecutorClient, w http.ResponseWriter, r *http.Request) {
	executeReply, err := client.ExecuteCmd(context.Background(), &gRPC.ExecuteMsg{
		Cmd: "delall",
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := DelAllResponse{
		ResponseMessage: utils.ResponseMessage{
			Message: "DelAll Success",
		},
		Data: utils.KeyVal{
			Key:   "",
			Value: executeReply.Value,
		},
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}
