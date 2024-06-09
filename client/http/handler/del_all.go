package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"tubes.sister/raft/client/http/utils"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

type DelAllResponse utils.KeyValResponse

func (gc *GRPCClient) DelAll(w http.ResponseWriter, r *http.Request) {
	executeReply, err := (*gc.client).ExecuteCmd(context.Background(), &gRPC.ExecuteMsg{
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
