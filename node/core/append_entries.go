package core

import (
	"strings"

	"tubes.sister/raft/node/data"
)

type AppendEntriesArgs struct {
	term         int             // leader’s term
	leaderId     data.Address    // so follower can redirect clients
	prevLogIndex int             // index of log entry immediately preceding new ones
	prevLogTerm  int             // term of prevLogIndex entry
	entries      []data.LogEntry // log entries to store (empty for heartbeat; may send more than one for efficiency)
	leaderCommit int             // leader’s commitIndex
}

type AppendEntriesReply struct {
	term    int  // currentTerm, for leader to update itself
	success bool // true if follower contained entry matching prevLogIndex and prevLogTerm
}

func (rn *RaftNode) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) error {
	// Rule 1 : Reply false if term < currentTerm (§5.1)
	if args.term < rn.Persistence.CurrentTerm {
		reply.success = false
		return nil
	}

	if args.term > rn.Persistence.CurrentTerm {
		rn.Persistence.CurrentTerm = args.term
		rn.Type = FOLLOWER
	}

	// Rule 2: Reply false if log doesn’t contain an
	// entry at prevLogIndex whose term matches prevLogTerm
	if len(rn.Persistence.Log) <= args.prevLogIndex || rn.Persistence.Log[args.prevLogIndex].Term != args.prevLogTerm {
		reply.success = false
		return nil
	}

	for i := args.prevLogIndex + 1; i <= args.prevLogIndex+len(args.entries); i++ {
		// Rule 3: If an existing entry conflicts with a new one (same index
		// but different terms), delete the existing entry and all that
		// follow it (§5.3)
		if len(rn.Persistence.Log)-1 >= i && rn.Persistence.Log[i].Term != args.term {
			rn.Persistence.Log = rn.Persistence.Log[0:i]
		}

		// Rule 4: Append any new entries not already in the log
		if len(rn.Persistence.Log) == i {
			rn.Persistence.Log = append(rn.Persistence.Log, args.entries[i-args.prevLogIndex-1])
		}
	}

	// Rule 5:  If leaderCommit > commitIndex, set commitIndex =
	// min(leaderCommit, index of last new entry)
	for i := rn.Volatile.CommitIndex + 1; i <= args.leaderCommit && i < len(rn.Persistence.Log); i++ {
		currentEntry := rn.Persistence.Log[i]

		command := currentEntry.Command
		// TODO change when key value mechanism is decided
		splitRes := strings.Split(currentEntry.Value, ",")
		key := splitRes[0]
		val := splitRes[1]

		switch command {
		case "APPEND":
			rn.Application.Append(key, val)
		case "SET":
			rn.Application.Set(key, val)
		case "DEL":
			rn.Application.Del(key)
		}

		rn.Volatile.CommitIndex = i
	}

	// TODO set leader address in node
	reply.term = rn.Persistence.CurrentTerm

	return nil
}
