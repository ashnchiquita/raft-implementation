package core

import (
	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/data"
)

// Updated RequestVote function
func (rn *RaftNode) RequestVote(args *gRPC.RequestVoteArgs) (*gRPC.RequestVoteReply, error) {
	// Node
	currTerm := rn.Persistence.CurrentTerm
	votedFor := rn.Persistence.VotedFor
	log := rn.Persistence.Log

	// Reply
	reply := &gRPC.RequestVoteReply{}

	// Rule 1 : Reply false if term < currentTerm (§5.1)
	if int(args.Term) < currTerm {
		reply.Term = int32(currTerm)
		reply.VoteGranted = false
		return reply, nil
	}

	// Reset election term if Term > currentTerm
	if int(args.Term) > currTerm {
		rn.Persistence.CurrentTerm = int(args.Term)
		rn.Volatile.LeaderAddress = data.Address{}
		rn.Volatile.Type = data.FOLLOWER
		if rn.Volatile.Type == data.CANDIDATE {
			rn.electionInterrupt <- HIGHER_TERM
			rn.Persistence.VotedFor = *data.NewZeroAddress()
		}
		// TODO: what to do if it's a leader?
	}

	// Rule 2 : If votedFor is null or candidateId, and candidate’s log is at least as up-to-date as receiver’s log, grant vote (§5.2, §5.4)
	candidateAddr := data.Address{IP: args.CandidateAddress.Ip, Port: int(args.CandidateAddress.Port)}
	if (votedFor.IsZeroAddress() ||
		candidateAddr.Equals(&votedFor)) &&
		isUpdated(int(args.LastLogIndex), int(args.LastLogTerm), log) {
		rn.Persistence.VotedFor = candidateAddr
		reply.VoteGranted = true
	} else {
		reply.VoteGranted = false
	}

	reply.Term = int32(currTerm)
	return reply, nil
}

func isUpdated(lastLogIndex int, lastLogTerm int, log []data.LogEntry) bool {
	if len(log) == 0 {
		return true
	}

	lastEntry := log[len(log)-1]
	if lastLogTerm != lastEntry.Term {
		return lastLogTerm > lastEntry.Term
	}
	return lastLogIndex >= len(log)
}
