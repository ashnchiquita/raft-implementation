package application

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"

	"tubes.sister/raft/client/http/utils"
)

type StrlenResponse struct {
	utils.ResponseMessage
	Data int `json:"data"`
}

func Strlen(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	resp := StrlenResponse{
		ResponseMessage: utils.ResponseMessage{
			Message: "success",
		},
		Data: len(key), // TODO: strlen key
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
