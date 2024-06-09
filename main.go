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
	programType   = "server"
	port          = flag.Int("port", 50051, "GRPC Server Port")
	serverAddr    = flag.String("server_address", "localhost:50051", "GRPC Server Address")
	clientPort    = flag.Int("client_port", 3000, "HTTP Client Port")
	clusterConfig = flag.String("clusterconfig", "./clusterconfig.json", "Initial cluster configuration")
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

	default: // server
		log.Println("Starting server...")

		addr := data.NewAddress("127.0.0.1", *port)

		raftNode := core.NewRaftNode(*addr)
		//! Remove these later
		// if *port == 5000 {
		// 	raftNode.InitializeAsLeader()
		// }
		// if *port == 5000 {
		// 	raftNode.Volatile.Type = data.LEADER
		// }
		// raftNode.Volatile.ClusterList = []data.ClusterData{
		// 	*data.NewClusterData(*data.NewAddress("localhost", 5000), 0, -1),
		// 	*data.NewClusterData(*data.NewAddress("localhost", 5001), 0, -1),
		// 	*data.NewClusterData(*data.NewAddress("localhost", 5002), 0, -1),
		// }

		raftNode.ReadClusterConfigFromFile(*clusterConfig)

		// raftNode.Volatile.LeaderAddress = *data.NewAddress("localhost", 5000)
		raftNode.InitializeServer()
	}
}

func parseFlags() {
	flag.Func("type", "'server' (default) or 'heartbeat' or 'terminal-client' or 'web-client'", func(flagValue string) error {
		if flagValue == "terminal-client" || flagValue == "web-client" || flagValue == "heartbeat" {
			programType = flagValue
		}

		return nil
	})

	flag.Parse()
}

// config [{"ip":"localhost","port":5000},{"ip":"localhost","port":5001}]
// config [{"ip":"127.0.0.1","port":5000},{"ip":"127.0.0.1","port":5001}]
