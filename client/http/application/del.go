package application

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"

	"tubes.sister/raft/client/http/utils"
)

type DelResponse utils.KeyValResponse

func Del(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	// TODO: del key

	dummy := utils.KeyVal{
		Key:   key,
		Value: "Dummy value (TODO: get from server)",
	}

	resp := GetResponse{
		ResponseMessage: utils.ResponseMessage{
			Message: "success",
		},
		Data: dummy,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
