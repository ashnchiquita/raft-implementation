package core

import (
	"context"
	"strings"

	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/data"
)

func (rn *RaftNode) AppendEntries(ctx context.Context, args *gRPC.AppendEntriesArgs) (*gRPC.AppendEntriesReply, error) {
	reply := &gRPC.AppendEntriesReply{Term: int32(rn.Persistence.CurrentTerm)}
	rn.setTimoutSafe(RandomizeElectionTimeout())

	// Rule 1 : Reply false if term < currentTerm (§5.1)
	if int(args.Term) < rn.Persistence.CurrentTerm {
		reply.Success = false
		return reply, nil
	}

	if int(args.Term) > rn.Persistence.CurrentTerm {
		rn.Persistence.CurrentTerm = int(args.Term)
		rn.Volatile.Type = data.FOLLOWER
	}

	// Rule 2: Reply false if log doesn’t contain an
	// entry at prevLogIndex whose term matches prevLogTerm
	if len(rn.Persistence.Log) <= int(args.PrevLogIndex) && (int(args.PrevLogIndex) != -1 && rn.Persistence.Log[args.PrevLogIndex].Term != int(args.PrevLogTerm)) {
		reply.Success = false
		return reply, nil
	}

	for i := int(args.PrevLogIndex) + 1; i <= int(args.PrevLogIndex)+len(args.Entries); i++ {
		// Rule 3: If an existing entry conflicts with a new one (same index
		// but different terms), delete the existing entry and all that
		// follow it (§5.3)
		if len(rn.Persistence.Log)-1 >= i && rn.Persistence.Log[i].Term != int(args.Term) {
			rn.Persistence.Log = rn.Persistence.Log[0:i]
		}

		// Rule 4: Append any new entries not already in the log
		if len(rn.Persistence.Log) == i {
			argsEntry := args.Entries[i-int(args.PrevLogIndex)-1]
			newEntry := data.LogEntry{Term: int(argsEntry.Term), Command: argsEntry.Command, Value: argsEntry.Value}
			rn.Persistence.Log = append(rn.Persistence.Log, newEntry)
			rn.Persistence.Serialize()
		}
	}

	// Rule 5:  If leaderCommit > commitIndex, set commitIndex =
	// min(leaderCommit, index of last new entry)
	for i := rn.Volatile.CommitIndex + 1; i <= int(args.LeaderCommit) && i < len(rn.Persistence.Log); i++ {
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

	rn.Volatile.LeaderAddress.IP = args.LeaderAddress.Ip
	rn.Volatile.LeaderAddress.Port = int(args.LeaderAddress.Port)

	reply.Success = true
	return reply, nil
}
