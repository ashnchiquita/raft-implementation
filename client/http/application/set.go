package application

import (
	"encoding/json"
	"net/http"

	"tubes.sister/raft/client/http/utils"
)

type SetRequest utils.KeyVal

func Set(w http.ResponseWriter, r *http.Request) {
	var setReq SetRequest
	json.NewDecoder(r.Body).Decode(&setReq)

	if setReq.Key == "" || setReq.Value == "" {
		msg := "key and value cannot be empty"
		utils.SendResponseMessage(w, msg, http.StatusBadRequest)
		return
	}

	// TODO: set key

	msg := "success"
	utils.SendResponseMessage(w, msg, http.StatusOK)
}
