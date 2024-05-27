package hello

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"tubes.sister/raft/server"
)

type HelloServerImpl struct {
	UnimplementedHelloServer
	raftNode *server.RaftNode
}

func (s *HelloServerImpl) SayHello(ctx context.Context, msg *HelloMsg) (*HelloRes, error) {
	log.Printf("Received: %v", msg.Name)
	return &HelloRes{Message: "Hello " + msg.Name}, nil
}

func NewHelloServerImpl(address server.Address) *HelloServerImpl {
	s := &HelloServerImpl{}
	s.raftNode = server.NewRaftNode(address)
	return s
}

func (s *HelloServerImpl) AddNode(ctx context.Context, req *AddNodeRequest) (*AddNodeResponse, error) {
	log.Printf("Received AddNode request for: %v", req.NewNodeAddress)

	// Pisah address jadi IP dan Port
	parts := strings.Split(req.NewNodeAddress, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid address format")
	}

	addr := server.Address{
		IP:   parts[0],
		Port: parts[1],
	}

	// Cek s.raftNode kosong
	if s.raftNode == nil {
		return nil, errors.New("raftNode is not initialized")
	}

	s.raftNode.AddNode(addr)

	return &AddNodeResponse{Success: true}, nil
}
