: ${MASTER_PORT:=5000}
: ${WEB_CLIENT_PORT:=3000}

go run main.go -type=web-client -server_address=localhost:$MASTER_PORT -client_port=$WEB_CLIENT_PORT
