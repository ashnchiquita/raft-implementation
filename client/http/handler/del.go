package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"tubes.sister/raft/client/http/utils"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

type DeleteResponse utils.KeyValResponse

// @Summary Delete value by key
// @ID delete-value-by-key
// @Tags         application
// @Produce      json
// @Param key path string true "Key to delete"
// @Success 200 {object} DeleteResponse
// @Failure 500 {object} utils.ResponseMessage
// @Router /app/{key} [delete]
func (gc *GRPCClient) Delete(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	executeReply, err := (*gc.client).ExecuteCmd(context.Background(), &gRPC.ExecuteMsg{
		Cmd:  "del",
		Vals: []string{key},
	})

	errMsg := "Failed to delete key-value pair"

	if err != nil {
		utils.SendResponseMessage(w, errMsg, http.StatusInternalServerError)
		return
	}

	resp := GetResponse{
		ResponseMessage: utils.ResponseMessage{
			Message: "Delete Success",
		},
		Data: utils.KeyVal{
			Key:   key,
			Value: executeReply.Value,
		},
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		utils.SendResponseMessage(w, errMsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}
