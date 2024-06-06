: ${MASTER_PORT:=5001}

go run main.go -type=heartbeat -port=$MASTER_PORT
