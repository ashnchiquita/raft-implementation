package core

import (
	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/data"
)

func (rn *RaftNode) RequestVote(args *gRPC.RequestVoteArgs, reply *gRPC.RequestVoteReply) error {

	// Rule 1 : Reply false if term < currentTerm (§5.1)
	if int(args.Term) < rn.Persistence.CurrentTerm {
		reply.Term = int32(rn.Persistence.CurrentTerm)
		reply.VoteGranted = false
		return nil
	}

	// Reset election term if Term > currentTerm
	if int(args.Term) > rn.Persistence.CurrentTerm {
		rn.Persistence.CurrentTerm = int(args.Term)
		rn.Volatile.LeaderAddress = data.Address{}
		rn.Volatile.Type = data.FOLLOWER
	}

	// Rule 2 : If votedFor is null or candidateId, and candidate’s log is at least as up-to-date as receiver’s log, grant vote (§5.2, §5.4)
	if (rn.Persistence.VotedFor == data.Address{} ||
		rn.Persistence.VotedFor == rn.Volatile.ClusterList[args.CandidateId].Address) &&
		isUpdated(int(args.LastLogIndex), int(args.LastLogTerm), rn.Persistence.Log) {
		rn.Persistence.VotedFor = rn.Volatile.ClusterList[args.CandidateId].Address
		reply.VoteGranted = true
	} else {
		reply.VoteGranted = false
	}

	reply.Term = int32(rn.Persistence.CurrentTerm)
	return nil
}

func isUpdated(lastLogIndex int, lastLogTerm int, log []data.LogEntry) bool {
	if lastLogIndex >= len(log) {
		return true
	}

	return lastLogTerm >= len(log)
}
