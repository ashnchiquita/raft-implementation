package utils

import (
	"encoding/json"
	"net/http"
)

func SendResponseMessage(w http.ResponseWriter, msg string, status int) {
	message := ResponseMessage{
		Message: msg,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(message)
}
