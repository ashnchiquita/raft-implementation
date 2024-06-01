package core

import (
	"context"
	"log"

	gRPC "tubes.sister/raft/gRPC/node/core"
)

func (n *RaftNode) SayHello(ctx context.Context, msg *gRPC.HelloMsg) (*gRPC.HelloRes, error) {
	log.Printf("Received: %v", msg.Name)
	return &gRPC.HelloRes{Message: "Hello " + msg.Name}, nil
}
