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
	"tubes.sister/raft/client/utils"
	gRPC "tubes.sister/raft/gRPC/node/core"
)

type TerminalClient struct {
	conn   *grpc.ClientConn
	client gRPC.CmdExecutorClient
}

func NewTerminalClient(clientPort int, serverAddr string) *TerminalClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(serverAddr, opts...)

	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	client := gRPC.NewCmdExecutorClient(conn)

	return &TerminalClient{conn: conn, client: client}
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
		"log":    1,
		"exit":   1,
		"config": 2,
	}

	commandLen, ok := validCommands[splitted[0]]
	if !ok {
		return fmt.Errorf("invalid command")
	}

	if len(splitted) != commandLen {
		return fmt.Errorf("invalid number of arguments")
	}

	if splitted[0] == "set" || splitted[0] == "append" {
		for _, arg := range splitted[1:] {
			if !utils.IsValidKeyOrValue(strings.TrimSpace(arg)) {
				return fmt.Errorf("%s command have to have valid key and value", splitted[0])
			}
		}
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

		executeReply, err := tc.client.ExecuteCmd(context.Background(), &gRPC.ExecuteMsg{
			Cmd:  splitted[0],
			Vals: splitted[1:],
		})
		if err != nil {
			fmt.Printf("Error executing command: %v\n", err)
		} else {
			fmt.Printf("ExecuteReply: %v\n", executeReply.Value)
		}
	}
}

func (tc *TerminalClient) Stop() {
	if tc.conn == nil {
		return
	}

	tc.conn.Close()
}
