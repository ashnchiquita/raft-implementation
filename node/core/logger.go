package core

import (
	"log"

	"tubes.sister/raft/node/data"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

func (rn *RaftNode) logf(format string, v ...interface{}) {
	nodetype := ""

	if rn.Volatile.Type == data.FOLLOWER {
		nodetype = "FOLLOWER"
	} else if rn.Volatile.Type == data.CANDIDATE {
		nodetype = "CANDIDATE"
	} else if rn.Volatile.Type == data.LEADER {
		nodetype = "LEADER"
	}
	coloredAddress := Red + "[" + rn.Address.String() + "]" + Cyan + "[" + nodetype + "]" + Reset

	log.Printf("%s "+format, append([]interface{}{coloredAddress}, v...)...)
}

func (rn *RaftNode) announcef(format string, v ...interface{}) {
	coloredAddress := Purple + format + Reset
	log.Printf("%s ", append([]interface{}{coloredAddress}, v...)...)
}
