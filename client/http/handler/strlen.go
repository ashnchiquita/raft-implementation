package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"tubes.sister/raft/client/http/utils"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

type StrlenResponse utils.KeyValResponse

func (gc *GRPCClient) Strlen(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	executeReply, err := (*gc.client).ExecuteCmd(context.Background(), &gRPC.ExecuteMsg{
		Cmd:  "strlen",
		Vals: []string{key},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := StrlenResponse{
		ResponseMessage: utils.ResponseMessage{
			Message: "Strlen Success",
		},
		Data: utils.KeyVal{
			Key:   key,
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
