package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"tubes.sister/raft/client/http/utils"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

type StrlenResponse struct {
	utils.ResponseMessage
	Data int `json:"data"`
}

// @Summary Get value string length by key
// @ID get-value-string-length-by-key
// @Tags         application
// @Produce      json
// @Param key path string true "Key to get"
// @Success 200 {object} StrlenResponse
// @Failure 500 {object} utils.ResponseMessage
// @Router /app/{key}/strlen [get]
func (gc *GRPCClient) Strlen(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	executeReply, err := (*gc.client).ExecuteCmd(context.Background(), &gRPC.ExecuteMsg{
		Cmd:  "strlen",
		Vals: []string{key},
	})

	errMsg := "Failed to get value length of key-value pair"

	if err != nil {
		utils.SendResponseMessage(w, errMsg, http.StatusInternalServerError)
		return
	}

	valLen, err := strconv.Atoi(executeReply.Value)
	if err != nil {
		utils.SendResponseMessage(w, errMsg, http.StatusInternalServerError)
		return
	}

	resp := StrlenResponse{
		ResponseMessage: utils.ResponseMessage{
			Message: "Strlen Success",
		},
		Data: valLen,
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		utils.SendResponseMessage(w, errMsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}
