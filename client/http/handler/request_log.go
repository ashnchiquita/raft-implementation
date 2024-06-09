package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"tubes.sister/raft/client/http/utils"
	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/data"
)

type RequestLogResponse struct {
	utils.ResponseMessage
	Data []data.LogEntry `json:"data"`
}

// @Summary Request cluster leader's log
// @ID request-log
// @Tags         cluster
// @Produce      json
// @Success 200 {object} RequestLogResponse
// @Failure 500 {object} utils.ResponseMessage
// @Router /cluster/log [get]
func (gc *GRPCClient) RequestLog(w http.ResponseWriter, r *http.Request) {
	executeReply, err := (*gc.client).ExecuteCmd(context.Background(), &gRPC.ExecuteMsg{
		Cmd:  "log",
	})

	errMsg := "Failed to request cluster leader's log"

	if err != nil {
		log.Println(errMsg + ": " + err.Error())
		utils.SendResponseMessage(w, errMsg, http.StatusInternalServerError)
		return
	}

	var data []data.LogEntry
	err = json.Unmarshal([]byte(executeReply.Value), &data)
	if err != nil {
		log.Println(errMsg + ": " + err.Error())
		utils.SendResponseMessage(w, errMsg, http.StatusInternalServerError)
		return
	}

	resp := RequestLogResponse{
		ResponseMessage: utils.ResponseMessage{
			Message: "success",
		},
		Data: data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
