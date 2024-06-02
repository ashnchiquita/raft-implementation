: ${MASTER_PORT:=5052}

go run main.go -type=server -port=$MASTER_PORT
