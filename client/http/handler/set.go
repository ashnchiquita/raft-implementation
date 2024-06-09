package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"tubes.sister/raft/client/http/utils"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

type SetResponse utils.KeyValResponse

// @Summary Set value to key
// @ID set-value-to-key
// @Tags         application
// @Accept       json
// @Produce      json
// @Param key body utils.KeyVal true "Key and value to set"
// @Success 200 {object} SetResponse
// @Router /app [put]
func (gc *GRPCClient) Set(w http.ResponseWriter, r *http.Request) {
	var kv utils.KeyVal
	err := json.NewDecoder(r.Body).Decode(&kv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	executeReply, err := (*gc.client).ExecuteCmd(context.Background(), &gRPC.ExecuteMsg{
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
