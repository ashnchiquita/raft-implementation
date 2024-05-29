package application

import (
	"encoding/json"
	"net/http"

	"tubes.sister/raft/client/http/utils"
)

type GetAllResponse utils.KeyValsResponse

func GetAll(w http.ResponseWriter, r *http.Request) {
	// TODO: get all keys

	dummy := []utils.KeyVal{
		{
			Key:   "Dummy key 1",
			Value: "Dummy value 1 (TODO: get from server)",
		},
		{
			Key:   "Dummy key 2",
			Value: "Dummy value 2 (TODO: get from server)",
		},
	}

	resp := GetAllResponse{
		ResponseMessage: utils.ResponseMessage{
			Message: "success",
		},
		Data: dummy,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
