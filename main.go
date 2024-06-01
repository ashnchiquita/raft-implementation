package main

import (
	"flag"
	"log"

	"tubes.sister/raft/client/http"
	"tubes.sister/raft/client/terminal"
)

var (
	programType = "server"
	port        = flag.Int("port", 50051, "GRPC Server Port")
	serverAddr  = flag.String("server_address", "localhost:50051", "GRPC Server Address")
	clientPort  = flag.Int("client_port", 3000, "HTTP Client Port")
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
		// TODO: start server
	}
}

func parseFlags() {
	flag.Func("type", "'server' (default) or 'terminal-client' or 'web-client'", func(flagValue string) error {
		if flagValue == "terminal-client" || flagValue == "web-client" {
			programType = flagValue
		}

		return nil
	})

	flag.Parse()
}
