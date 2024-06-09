package core

import (
	"context"
	"log"

	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/data"
)

// Updated RequestVote function
func (rn *RaftNode) RequestVote(ctx context.Context, args *gRPC.RequestVoteArgs) (*gRPC.RequestVoteReply, error) {
	// Node
	log := Green + "RequestVote() >> " + Reset
	currTerm := rn.Persistence.CurrentTerm
	votedFor := rn.Persistence.VotedFor
	currLog := rn.Persistence.Log
	candidateAddr := data.Address{IP: args.CandidateAddress.Ip, Port: int(args.CandidateAddress.Port)}
	rn.logf(log+"RPC received from %v; with Term: %d; Last Log Index: %d; Last log term: %d", args.CandidateAddress.Port, args.Term, args.LastLogIndex, args.LastLogTerm)

	// Reply
	reply := &gRPC.RequestVoteReply{SameCluster: false}

	// Will tell candidate if the candidate is in the same cluster or not
	// If candidate is not in a joint consesnsus state, the candidate will
	// stop it's process when it hits a follower in a different cluster
	for _, clusterData := range rn.Volatile.ClusterList {
		if clusterData.Address.Equals(&candidateAddr) {
			reply.SameCluster = true
			break
		}
	}

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
		rn.Volatile.LeaderAddress = *data.NewZeroAddress()
		rn.Persistence.VotedFor = *data.NewZeroAddress()

		if rn.Volatile.Type == data.CANDIDATE {
			rn.logf(log + "Term is higher than current term. Going to follower state from candidate.")
			rn.electionInterrupt <- HIGHER_TERM
			rn.cleanupCandidateState()
		} else if rn.Volatile.Type == data.LEADER {
			rn.logf(log + "Term is higher than current term. Going to follower state from leader.")
		}
		rn.setAsFollower()
	}

	// Rule 2 : If votedFor is null or candidateId, and candidate’s log is at least as up-to-date as receiver’s log, grant vote (§5.2, §5.4)
	if (votedFor.IsZeroAddress() ||
		candidateAddr.Equals(&votedFor)) &&
		isUpdated(int(args.LastLogIndex), int(args.LastLogTerm), currLog) {
		rn.Persistence.VotedFor = candidateAddr
		reply.VoteGranted = true
		rn.logf(log + "Vote granted")
	} else {
		reply.VoteGranted = false
		if !(isUpdated(int(args.LastLogIndex), int(args.LastLogTerm), currLog)) {
			rn.logf(log + "Vote not granted because candidate's log is not up-to-date")
		} else {
			rn.logf(log + "Vote not granted because votedFor is not null or candidateId")
		}
	}

	reply.Term = int32(currTerm)

	// Persist after accepting vote request
	rn.Persistence.Serialize()

	return reply, nil
}

func isUpdated(lastLogIndex int, lastLogTerm int, currLog []data.LogEntry) bool {
	if len(currLog) == 0 {
		return true
	}

	lastEntry := currLog[len(currLog)-1]
	log.Printf(Green+"RequestVote() >> "+Reset+"Last entry index: %d; Last Entry term: %d", len(currLog)-1, lastEntry.Term)
	if lastLogTerm != lastEntry.Term {
		return lastLogTerm >= lastEntry.Term
	}
	return lastLogIndex >= len(currLog)-1
}
