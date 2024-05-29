package application

import (
	"net/http"

	"tubes.sister/raft/client/http/utils"
)

func DelAll(w http.ResponseWriter, r *http.Request) {
	// TODO: delete all key

	msg := "success"
	utils.SendResponseMessage(w, msg, http.StatusOK)
}
