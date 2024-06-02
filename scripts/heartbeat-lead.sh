: ${MASTER_PORT:=5051}
: ${MASTER_FOLLOWER:=5052}

go run main.go -type=heartbeat-lead -port=$MASTER_PORT -follower=$MASTER_FOLLOWER
