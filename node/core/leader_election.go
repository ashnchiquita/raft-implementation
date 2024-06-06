package core

import (
	"context"
	"fmt"
	"log"
	"math"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/data"
)

type voteResult struct {
	term        int
	voteGranted bool
}

type voteRequest struct {
	candidateAddress data.Address
	lastTerm         int
	lastIndex        int
}

func (rn *RaftNode) cleanupCandidateState() {
	rn.Volatile.ResetVotes()
	rn.Persistence.VotedFor = *data.NewZeroAddress()
}

func (rn *RaftNode) setAsFollower() {
	rn.Volatile.Type = data.FOLLOWER
	rn.resetTimeout()
}

func (rn *RaftNode) setAsCandidate() {
	rn.Volatile.Type = data.CANDIDATE
	rn.resetTimeout()
}

func (rn *RaftNode) startElection() {
	c := make(chan bool)

	go rn.election(c)
	for !<-c {
		go rn.election(c)
	}
}

// Will convert the current node to a candidate (+ reset its timeout) and start the election process
func (rn *RaftNode) election(c chan bool) {
	votes := make(chan voteResult)

	rn.setAsCandidate()

	rn.Persistence.CurrentTerm++
	rn.Persistence.VotedFor = rn.Address
	rn.Volatile.AddVote(rn.Address)
	lastTerm := 0
	if len(rn.Persistence.Log) > 0 {
		lastTerm = rn.Persistence.Log[len(rn.Persistence.Log)-1].Term
	}

	voteReq := voteRequest{
		candidateAddress: rn.Address,
		lastTerm:         lastTerm,
		lastIndex:        len(rn.Persistence.Log) - 1,
	}

	for _, node := range rn.Volatile.ClusterList {
		if node.Address == rn.Address {
			continue
		}

		go rn.requestVote(&voteReq, &node, votes)
	}

	for {
		// accept reply until interrupted
		select {
		case val := <-rn.electionInterrupt:
			rn.cleanupCandidateState()

			switch val {
			case ELECTION_TIMEOUT:
				// restart election (auto restart because triggered by timeout)
				c <- false
			case HIGHER_TERM:
				rn.setAsFollower()
				c <- true
			}

			return
		case vote := <-votes:
			if rn.Volatile.Type == data.CANDIDATE && vote.term == rn.Persistence.CurrentTerm && vote.voteGranted {
				rn.Volatile.AddVote(rn.Address)
				if rn.Volatile.GetVotesCount() >= int(math.Ceil(float64(len(rn.Volatile.ClusterList))/2+1)) {
					rn.cleanupCandidateState()
					rn.InitializeAsLeader() // will also reset timeout
					return
				}
			} else if vote.term > rn.Persistence.CurrentTerm {
				rn.cleanupCandidateState()
				rn.Persistence.CurrentTerm = vote.term
				rn.setAsFollower() // will also reset timeout
				return
			}
		}
	}
}

func (rn *RaftNode) requestVote(voteReq *voteRequest, node *data.ClusterData, votes chan voteResult) {
	log.Printf("Requesting vote to to %v", node.Address)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", node.Address.IP, node.Address.Port), opts...)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	client := gRPC.NewRaftNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), rn.timeout.Value)

	defer cancel()

	reply, err := client.RequestVote(ctx, &gRPC.RequestVoteArgs{
		Term:             int32(rn.Persistence.CurrentTerm),
		CandidateAddress: &gRPC.RequestVoteArgs_CandidateAddress{Ip: voteReq.candidateAddress.IP, Port: int32(voteReq.candidateAddress.Port)},
		LastLogIndex:     int32(voteReq.lastIndex),
		LastLogTerm:      int32(voteReq.lastTerm),
	})

	if err != nil {
		log.Printf("Error requesting vote to %v: %v", node.Address, err)
		return
	}

	votes <- voteResult{term: rn.Persistence.CurrentTerm, voteGranted: reply.VoteGranted}
}
