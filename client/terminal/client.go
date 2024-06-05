package terminal

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

type TerminalClient struct {
	conn *grpc.ClientConn
}

func NewTerminalClient(clientPort int, serverAddr string) *TerminalClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(serverAddr, opts...)

	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	return &TerminalClient{conn: conn}
}

func validateInput(splitted []string) error {
	validCommands := map[string]int{
		"ping":   1,
		"get":    2,
		"set":    3,
		"strlen": 2,
		"del":    2,
		"append": 3,
		"getall": 1,
		"delall": 1,
		"exit":   1,
	}

	commandLen, ok := validCommands[splitted[0]]
	if !ok {
		return fmt.Errorf("invalid command")
	}

	if len(splitted) != commandLen {
		return fmt.Errorf("invalid number of arguments")
	}

	return nil
}

func (tc *TerminalClient) Start() {
	scanner := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Command: ")
		input, inputErr := scanner.ReadString('\n')

		if inputErr != nil {
			fmt.Println(inputErr)
			break
		}

		input = strings.TrimSpace(input)
		splitted := strings.Split(input, " ")

		if err := validateInput(splitted); err != nil {
			fmt.Println(err)
			continue
		}

		if splitted[0] == "exit" {
			fmt.Println("Exiting terminal client...")
			break
		}

		// TODO: Determine whether creating a new client for each command input is better
		// TODO: than creating a single client as attribute to TerminalClient
		client := gRPC.NewCmdExecutorClient(tc.conn)

		// TODO: implement the command
		fmt.Println("command", input, "called")
		executeReply, err := client.ExecuteCmd(context.Background(), &gRPC.ExecuteMsg{
			Cmd:  "set",
			Vals: splitted[1:],
		})
		if err != nil {
			log.Printf("Error executing command: %v", err)
		}

		log.Printf("ExecuteReply: %v", executeReply.Success)
	}
}

func (tc *TerminalClient) Stop() {
	if tc.conn == nil {
		return
	}

	tc.conn.Close()
}
