package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"tubes.sister/raft/client/http/utils"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

type GetResponse utils.KeyValResponse

// @Summary Get value by key
// @ID get-value-by-key
// @Tags         application
// @Produce      json
// @Param key path string true "Key to get"
// @Success 200 {object} GetResponse
// @Router /app/{key} [get]
func (gc *GRPCClient) Get(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	executeReply, err := (*gc.client).ExecuteCmd(context.Background(), &gRPC.ExecuteMsg{
		Cmd:  "get",
		Vals: []string{key},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := GetResponse{
		ResponseMessage: utils.ResponseMessage{
			Message: "Get Success",
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
