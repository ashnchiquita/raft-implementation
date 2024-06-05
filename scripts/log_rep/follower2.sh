: ${FOLLOWER_PORT:=5002}

go run main.go -type=server -port=$FOLLOWER_PORT