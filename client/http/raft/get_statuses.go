package raft

import (
	"encoding/json"
	"net/http"

	"tubes.sister/raft/client/http/utils"
	"tubes.sister/raft/server"
)

// TODO: adjust type
type NodeStatus struct {
	Address server.Address `json:"address"`
	Alive   bool           `json:"alive"`
}

type GetStatusesResponse struct {
	utils.ResponseMessage
	Statuses []NodeStatus `json:"statuses"`
}

func GetStatuses(w http.ResponseWriter, r *http.Request) {
	// TODO: retrieve status from server

	dummy := []NodeStatus{
		{
			Address: server.Address{
				IP:   "127.0.0.1 (TODO: get from server)",
				Port: "8080",
			},
			Alive: true,
		},
		{
			Address: server.Address{
				IP:   "127.0.0.1 (TODO: get from server)",
				Port: "8081",
			},
			Alive: false,
		},
	}

	resp := GetStatusesResponse{
		ResponseMessage: utils.ResponseMessage{
			Message: "success",
		},
		Statuses: dummy,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
