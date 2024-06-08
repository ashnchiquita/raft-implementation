: ${FOLLOWER_PREFIX:=500}

go run main.go -type=server -port=$FOLLOWER_PREFIX$1
