package handler

import gRPC "tubes.sister/raft/gRPC/node/core"

type GRPCClient struct {
	client *gRPC.CmdExecutorClient
}

func NewGRPCClient(client *gRPC.CmdExecutorClient) *GRPCClient {
	return &GRPCClient{
		client: client,
	}
}
