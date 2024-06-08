package core

import (
	"context"
	"log"

	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/data"
)

// Updated RequestVote function
func (rn *RaftNode) RequestVote(ctx context.Context, args *gRPC.RequestVoteArgs) (*gRPC.RequestVoteReply, error) {
	log.Printf("RequestVote() >> RPC received from %v\n", args.CandidateAddress.Port)
	// Node
	currTerm := rn.Persistence.CurrentTerm
	votedFor := rn.Persistence.VotedFor
	currLog := rn.Persistence.Log

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
		currTerm = int(args.Term)
		rn.Volatile.LeaderAddress = data.Address{}

		if rn.Volatile.Type == data.CANDIDATE {
			log.Println("RequestVote() >> Term is higher than current term. Going to follower state from candidate.")
			rn.electionInterrupt <- HIGHER_TERM
			rn.cleanupCandidateState()
		} else if rn.Volatile.Type == data.LEADER {
			log.Println("RequestVote() >> Term is higher than current term. Going to follower state from leader.")
		}
		rn.setAsFollower()
	}

	// Rule 2 : If votedFor is null or candidateId, and candidate’s log is at least as up-to-date as receiver’s log, grant vote (§5.2, §5.4)
	candidateAddr := data.Address{IP: args.CandidateAddress.Ip, Port: int(args.CandidateAddress.Port)}
	if (votedFor.IsZeroAddress() ||
		candidateAddr.Equals(&votedFor)) &&
		isUpdated(int(args.LastLogIndex), int(args.LastLogTerm), currLog) {
		rn.Persistence.VotedFor = candidateAddr
		reply.VoteGranted = true
		log.Println("RequestVote() >> Vote granted")
	} else {
		reply.VoteGranted = false
		if !candidateAddr.Equals(&votedFor) {
			log.Println("RequestVote() >> Vote not granted because votedFor is not null or candidateId")
		} else {
			log.Println("RequestVote() >> Vote not granted because candidate's log is not up-to-date")
		}
	}

	reply.Term = int32(currTerm)
	return reply, nil
}

func isUpdated(lastLogIndex int, lastLogTerm int, currLog []data.LogEntry) bool {
	if len(currLog) == 0 {
		return true
	}

	lastEntry := currLog[len(currLog)-1]
	if lastLogTerm != lastEntry.Term {
		return lastLogTerm > lastEntry.Term
	}
	return lastLogIndex >= len(currLog)
}
