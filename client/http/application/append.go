package application

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"tubes.sister/raft/client/http/utils"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

type AppendResponse utils.KeyValResponse

func Append(client gRPC.CmdExecutorClient, w http.ResponseWriter, r *http.Request) {
	var req utils.KeyVal
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Append the new value
	appendReply, err := client.ExecuteCmd(context.Background(), &gRPC.ExecuteMsg{
		Cmd:  "append",
		Vals: []string{req.Key, req.Value},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := AppendResponse{
		ResponseMessage: utils.ResponseMessage{
			Message: "Append Success",
		},
		Data: utils.KeyVal{
			Key:   req.Key,
			Value: appendReply.Value,
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
