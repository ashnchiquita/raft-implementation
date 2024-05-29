package raft

import (
	"encoding/json"
	"net/http"

	"tubes.sister/raft/client/http/utils"
)

type PingResponse struct {
	utils.ResponseMessage
	Data string `json:"data"`
}

func Ping(w http.ResponseWriter, r *http.Request) {
	// TODO: ping server

	resp := PingResponse{
		ResponseMessage: utils.ResponseMessage{
			Message: "success",
		},
		Data: "PONG",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
