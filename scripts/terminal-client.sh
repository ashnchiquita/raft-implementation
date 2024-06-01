: ${MASTER_PORT:=5000}
: ${TERMINAL_CLIENT_PORT:=3001}

go run main.go -type=terminal-client -server_address=localhost:$MASTER_PORT -client_port=$TERMINAL_CLIENT_PORT
