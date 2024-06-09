package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"tubes.sister/raft/client/http/utils"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

type AppendRequest utils.KeyVal
type AppendResponse utils.KeyValResponse

// @Summary Append value to key
// @ID append-value-to-key
// @Tags         application
// @Accept       json
// @Produce      json
// @Param key body AppendRequest true "Key and value to append"
// @Success 200 {object} AppendResponse
// @Failure 400 {object} utils.ResponseMessage
// @Failure 500 {object} utils.ResponseMessage
// @Router /app [patch]
func (gc *GRPCClient) Append(w http.ResponseWriter, r *http.Request) {
	var appendReq AppendRequest
	json.NewDecoder(r.Body).Decode(&appendReq)

	if appendReq.Key == "" || appendReq.Value == "" {
		errMsg := "Key and value cannot be empty"
		utils.SendResponseMessage(w, errMsg, http.StatusBadRequest)
		return
	}

	errMsg := "Failed to append value to key"

	// Append the new value
	appendReply, err := (*gc.client).ExecuteCmd(context.Background(), &gRPC.ExecuteMsg{
		Cmd:  "append",
		Vals: []string{appendReq.Key, appendReq.Value},
	})
	if err != nil {
		utils.SendResponseMessage(w, errMsg, http.StatusInternalServerError)
		return
	}

	resp := AppendResponse{
		ResponseMessage: utils.ResponseMessage{
			Message: "Append Success",
		},
		Data: utils.KeyVal{
			Key:   appendReq.Key,
			Value: appendReply.Value,
		},
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		utils.SendResponseMessage(w, errMsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}
