: ${MASTER_PORT:=5000}

go run main.go -type=server -port=$MASTER_PORT
