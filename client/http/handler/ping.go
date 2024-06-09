package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"tubes.sister/raft/client/http/utils"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

type PingResponse struct {
	utils.ResponseMessage
	Data string `json:"data"`
}

// @Summary Ping cluster
// @ID ping-cluster
// @Tags         cluster
// @Produce      json
// @Success 200 {object} PingResponse
// @Failure 500 {object} utils.ResponseMessage
// @Router /cluster/ping [get]
func (gc *GRPCClient) Ping(w http.ResponseWriter, r *http.Request) {
	executeReply, err := (*gc.client).ExecuteCmd(context.Background(), &gRPC.ExecuteMsg{
		Cmd:  "ping",
	})

	errMsg := "Failed to ping cluster"

	if err != nil {
		log.Println(errMsg + ": " + err.Error())
		utils.SendResponseMessage(w, errMsg, http.StatusInternalServerError)
		return
	}

	resp := PingResponse{
		ResponseMessage: utils.ResponseMessage{
			Message: "success",
		},
		Data: executeReply.Value,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
