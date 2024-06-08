package application

import (
	"context"
	"encoding/json"
	"net/http"

	"tubes.sister/raft/client/http/utils"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

type SetResponse utils.KeyValResponse

func Set(client gRPC.CmdExecutorClient, w http.ResponseWriter, r *http.Request) {
	var kv utils.KeyVal
	err := json.NewDecoder(r.Body).Decode(&kv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	executeReply, err := client.ExecuteCmd(context.Background(), &gRPC.ExecuteMsg{
		Cmd:  "set",
		Vals: []string{kv.Key, kv.Value},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := SetResponse{
		ResponseMessage: utils.ResponseMessage{
			Message: "Set Success",
		},
		Data: utils.KeyVal{
			Key:   kv.Key,
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
