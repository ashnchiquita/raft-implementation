package main

import (
	"flag"
	"log"

	"tubes.sister/raft/client/http"

	"tubes.sister/raft/client/terminal"
	"tubes.sister/raft/node/core"
	"tubes.sister/raft/node/data"
)

var (
	programType  = "server"
	port         = flag.Int("port", 50051, "GRPC Server Port")
	followerPort = flag.Int("follower", 50052, "GRPC Follower Port for heartbeat testing")
	serverAddr   = flag.String("server_address", "localhost:50051", "GRPC Server Address")
	clientPort   = flag.Int("client_port", 3000, "HTTP Client Port")
)

func main() {
	parseFlags()

	switch programType {
	case "terminal-client":
		log.Println("Starting terminal client...")
		terminalClient := terminal.NewTerminalClient(*clientPort, *serverAddr)
		terminalClient.Start()
		defer terminalClient.Stop()

	case "web-client":
		log.Println("Starting web client...")
		webClient := http.NewHTTPClient(*clientPort, *serverAddr)
		webClient.Start()
		defer webClient.Stop()

	case "heartbeat-lead":
		log.Println("Starting heartbeat leader...")
		addr := data.NewAddress("localhost", *port)
		app := application.NewApplication()

		raftNode := core.NewRaftNode(*addr, *app)
		raftNode.Volatile.ClusterList = append(raftNode.Volatile.ClusterList, *data.NewAddress("localhost", *followerPort))
		raftNode.Volatile.Type = data.LEADER
		raftNode.InitializeServer()

	default: // server
		log.Println("Starting server...")

		addr := data.NewAddress("localhost", *port)

		raftNode := core.NewRaftNode(*addr)
		//! Remove these later
		if *port == 5000 {
			raftNode.Volatile.Type = data.LEADER
		}
		raftNode.Volatile.ClusterList = []data.ClusterData{
			{Address: *data.NewAddress("localhost", 5000), MatchIndex: -1, NextIndex: 0},
			{Address: *data.NewAddress("localhost", 5001), MatchIndex: -1, NextIndex: 0},
			{Address: *data.NewAddress("localhost", 5002), MatchIndex: -1, NextIndex: 0},
		}
		raftNode.InitializeServer()
	}
}

func parseFlags() {
	flag.Func("type", "'server' (default) or 'heartbeat-lead' or 'terminal-client' or 'web-client'", func(flagValue string) error {
		if flagValue == "terminal-client" || flagValue == "web-client" || flagValue == "heartbeat-lead" {
			programType = flagValue
		}

		return nil
	})

	flag.Parse()
}
