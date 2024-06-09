package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"tubes.sister/raft/client/http/utils"
	clientUtils "tubes.sister/raft/client/utils"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

type AppendRequest utils.KeyVal

// @Summary Append value to key
// @ID append-value-to-key
// @Tags         application
// @Accept       json
// @Produce      json
// @Param key body AppendRequest true "Key and value to append"
// @Success 200 {object} utils.ResponseMessage
// @Failure 400 {object} utils.ResponseMessage
// @Failure 500 {object} utils.ResponseMessage
// @Router /app [patch]
func (gc *GRPCClient) Append(w http.ResponseWriter, r *http.Request) {
	var appendReq AppendRequest
	json.NewDecoder(r.Body).Decode(&appendReq)

	if !clientUtils.IsValidPair(appendReq.Key, appendReq.Value) {
		errMsg := "Key and value cannot be empty, contain spaces, or contain commas"
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
		log.Println(errMsg + ": " + err.Error())
		utils.SendResponseMessage(w, errMsg, http.StatusInternalServerError)
		return
	}

	if !appendReply.Success {
		log.Println(errMsg + ": " + appendReply.Value)
		utils.SendResponseMessage(w, errMsg, http.StatusInternalServerError)
		return
	}

	msg := "Append Success"
	utils.SendResponseMessage(w, msg, http.StatusOK)
}
