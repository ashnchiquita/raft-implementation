package rpc

import (
	"tubes.sister/raft/server"
)

type RequestVoteArgs struct {
	term         int
	candidateId  server.Address
	lastLogIndex int
	lastLogTerm  int
}

type RequestVoteReply struct {
	term        int
	voteGranted bool
}

// Testing requirements
type TestNode struct {
	address           server.Address
	nodeType          server.NodeType
	log               []string
	electionTerm      int
	clusterAddrList   []server.Address
	clusterLeaderAddr server.Address
}

func (rn *TestNode) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) error {
	// Mutex lock
	// rn.mu.Lock()
	// defer rn.mu.Unlock()

	// Rule 1 : Reply false if term < currentTerm (§5.1)
	if args.term < rn.electionTerm {
		reply.term = rn.electionTerm
		reply.voteGranted = false
		return nil
	}

	// Reset election term if Term > currentTerm
	if args.term > rn.electionTerm {
		rn.electionTerm = args.term
		rn.clusterLeaderAddr = server.Address{}
		rn.nodeType = server.FOLLOWER
	}

	// Rule 2 : If votedFor is null or candidateId, and candidate’s log is at least as up-to-date as receiver’s log, grant vote (§5.2, §5.4)
	if (rn.clusterLeaderAddr == server.Address{} || rn.clusterLeaderAddr == args.candidateId) && isUpdated(args.lastLogIndex, args.lastLogTerm, rn.log) {
		rn.clusterLeaderAddr = args.candidateId
		reply.voteGranted = true
	} else {
		reply.voteGranted = false
	}

	reply.term = rn.electionTerm
	return nil
}

func isUpdated(lastLogIndex int, lastLogTerm int, log []string) bool {
	if lastLogIndex >= len(log) {
		return true
	}
	return lastLogTerm >= len(log)
}
