package hello

import (
	"context"
	"log"
)

type HelloServerImpl struct {
	UnimplementedHelloServer
}

func (s *HelloServerImpl) SayHello(ctx context.Context, msg *HelloMsg) (*HelloRes, error) {
	log.Printf("Received: %v", msg.Name)
	return &HelloRes{Message: "Hello " + msg.Name}, nil
}
