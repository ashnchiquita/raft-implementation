package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"tubes.sister/raft/client/http/utils"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

type GetAllResponse struct {
	utils.ResponseMessage
	Data []utils.KeyVal `json:"data"`
}

// @Summary Get all key-value pairs
// @ID get-all
// @Tags         application
// @Produce      json
// @Success 200 {object} GetAllResponse
// @Router /app [get]
func (gc *GRPCClient) GetAll(w http.ResponseWriter, r *http.Request) {
	executeReply, err := (*gc.client).ExecuteCmd(context.Background(), &gRPC.ExecuteMsg{
		Cmd: "getall",
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data []utils.KeyVal
	err = json.Unmarshal([]byte(executeReply.Value), &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := GetAllResponse{
		ResponseMessage: utils.ResponseMessage{
			Message: "GetAll Success",
		},
		Data: data,
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}
