package application

import (
	"encoding/json"
	"net/http"

	"tubes.sister/raft/client/http/utils"
)

type AppendRequest utils.KeyVal

func Append(w http.ResponseWriter, r *http.Request) {
	var appendReq AppendRequest
	json.NewDecoder(r.Body).Decode(&appendReq)

	if appendReq.Key == "" || appendReq.Value == "" {
		msg := "key and value cannot be empty"
		utils.SendResponseMessage(w, msg, http.StatusBadRequest)
		return
	}

	// TODO: append key

	msg := "success"
	utils.SendResponseMessage(w, msg, http.StatusOK)
}
