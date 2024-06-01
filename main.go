package main

import (
	"flag"
	"log"

	"tubes.sister/raft/client/http"
	"tubes.sister/raft/client/terminal"
)

var (
	programType = "server"
	port       	= flag.Int("port", 50051, "GRPC Server Port")
	// serverAddr 	= flag.String("server_address", "localhost:50051", "GRPC Server Address")
	clientPort	= flag.Int("client_port", 3000, "HTTP Client Port")
)

func main() {
	parseFlags()

	switch programType {
		case "terminal-client":
			log.Println("Starting terminal client...")
			terminal.StartTerminalClient(*clientPort, *port)

		case "web-client":
			log.Println("Starting web client...")
			http.StartHTTPClient(*clientPort, *port)

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
