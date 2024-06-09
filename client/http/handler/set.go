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

type SetRequest utils.KeyVal
type SetResponse utils.KeyValResponse

// @Summary Set value to key
// @ID set-value-to-key
// @Tags         application
// @Accept       json
// @Produce      json
// @Param key body SetRequest true "Key and value to set"
// @Success 200 {object} SetResponse
// @Failure 400 {object} utils.ResponseMessage
// @Failure 500 {object} utils.ResponseMessage
// @Router /app [put]
func (gc *GRPCClient) Set(w http.ResponseWriter, r *http.Request) {
	var setReq SetRequest
	json.NewDecoder(r.Body).Decode(&setReq)

	if !clientUtils.IsValidPair(setReq.Key, setReq.Value) {
		errMsg := "Key and value cannot be empty, contain spaces, or contain commas"
		utils.SendResponseMessage(w, errMsg, http.StatusBadRequest)
		return
	}

	executeReply, err := (*gc.client).ExecuteCmd(context.Background(), &gRPC.ExecuteMsg{
		Cmd:  "set",
		Vals: []string{setReq.Key, setReq.Value},
	})

	errMsg := "Failed to set key-value pair"

	if err != nil {
		log.Println(errMsg + ": " + err.Error())
		utils.SendResponseMessage(w, errMsg, http.StatusInternalServerError)
		return
	}

	resp := SetResponse{
		ResponseMessage: utils.ResponseMessage{
			Message: "Set Success",
		},
		Data: utils.KeyVal{
			Key:   setReq.Key,
			Value: executeReply.Value,
		},
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Println(errMsg + ": " + err.Error())
		utils.SendResponseMessage(w, errMsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}
