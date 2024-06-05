: ${FOLLOWER_PORT:=5001}

go run main.go -type=server -port=$FOLLOWER_PORT
