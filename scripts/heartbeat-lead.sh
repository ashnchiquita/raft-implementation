: ${MASTER_PORT:=5000}

go run main.go -type=heartbeat-lead -port=$MASTER_PORT
