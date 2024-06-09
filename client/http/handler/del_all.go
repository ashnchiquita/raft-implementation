package handler

import (
	"context"
	"net/http"

	"tubes.sister/raft/client/http/utils"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

// @Summary Delete all key-value pairs
// @ID delete-all
// @Tags         application
// @Produce      json
// @Success 200 {object} utils.ResponseMessage
// @Failure 500 {object} utils.ResponseMessage
// @Router /app [delete]
func (gc *GRPCClient) DelAll(w http.ResponseWriter, r *http.Request) {
	_, err := (*gc.client).ExecuteCmd(context.Background(), &gRPC.ExecuteMsg{
		Cmd: "delall",
	})

	if err != nil {
		errMsg := "Failed to delete all key-value pairs"
		utils.SendResponseMessage(w, errMsg, http.StatusInternalServerError)
		return
	}

	msg := "DelAll Success"
	utils.SendResponseMessage(w, msg, http.StatusOK)
}
