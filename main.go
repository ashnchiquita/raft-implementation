package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"tubes.sister/raft/hello"
)

var (
	server     = flag.Bool("is_server", true, "GRPC Server or Client")
	port       = flag.Int("port", 50051, "GRPC Server Port")
	serverAddr = flag.String("server_address", "localhost:50051", "GRPC Server Address")
)

func main() {
	flag.Parse()

	if *server {
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))

		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		var opts []grpc.ServerOption
		grpcServer := grpc.NewServer(opts...)
		hello.RegisterHelloServer(grpcServer, &hello.HelloServerImpl{})

		log.Printf("Server started at port %d", *port)
		grpcServer.Serve(lis)
	} else {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

		conn, err := grpc.Dial(*serverAddr, opts...)
		if err != nil {
			log.Fatalf("Failed to dial server: %v", err)
		}

		defer conn.Close()

		client := hello.NewHelloClient(conn)
		response, err := client.SayHello(context.Background(), &hello.HelloMsg{Name: "James"})
		if err != nil {
			log.Fatalf("Failed to call SayHello: %v", err)
		}

		log.Print(response.Message)
	}
}
