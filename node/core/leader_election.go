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
	term		int
	voteGranted	bool
}

type voteRequest struct {
	candidateAddress data.Address
	lastTerm 			 int
	lastIndex			 int
}

func (rn *RaftNode) backToFollower() {
	rn.Volatile.Type = data.FOLLOWER
	rn.Persistence.VotedFor = *data.NewZeroAddress()
	rn.Volatile.ResetVotes()
	rn.setTimeout()
}

func (rn *RaftNode) startElection() {
	votes := make(chan voteResult)

	rn.Volatile.Type = data.CANDIDATE
	rn.Volatile.ResetVotes()
	rn.setTimeout()
	rn.Persistence.CurrentTerm++
	rn.Persistence.VotedFor = rn.Address
	rn.Volatile.AddVote(rn.Address)
	lastTerm := 0
	if len(rn.Persistence.Log) > 0 {
		lastTerm = rn.Persistence.Log[len(rn.Persistence.Log)-1].Term
	}

	rn.Persistence.Serialize()

	voteReq := voteRequest{
		candidateAddress: rn.Address,
		lastTerm: lastTerm,
		lastIndex: len(rn.Persistence.Log) - 1,
	}

	for _, node := range rn.Volatile.ClusterList {
		if node.Address == rn.Address {
			continue
		}

		go rn.requestVote(&voteReq, &node, votes)
	}

	for {
		select {
			// case <-rn.timeoutElection:
				// TODO: restart election
				// return
			case vote := <-votes:
				if rn.Volatile.Type == data.CANDIDATE && vote.term == rn.Persistence.CurrentTerm && vote.voteGranted {
					rn.Volatile.AddVote(rn.Address)
					if rn.Volatile.GetVotesCount() >= int(math.Ceil(float64(len(rn.Volatile.ClusterList))/2 + 1)) {
						rn.InitializeAsLeader()
						return
					}
				} else if vote.term > rn.Persistence.CurrentTerm {
					rn.backToFollower()
					rn.Persistence.CurrentTerm = vote.term

					// TODO: reset election timeout
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

	reply, err := client.RequestVote(context.Background(), &gRPC.RequestVoteArgs{
		Term:         int32(rn.Persistence.CurrentTerm),
		CandidateAddress:  &gRPC.RequestVoteArgs_CandidateAddress{Ip: voteReq.candidateAddress.IP, Port: int32(voteReq.candidateAddress.Port)},
		LastLogIndex: int32(voteReq.lastIndex),
		LastLogTerm:  int32(voteReq.lastTerm),
	})

	if err != nil {
		log.Printf("Error requesting vote to %v: %v", node.Address, err)
		// TODO: Handle error? should it considered as vote not granted? secara tidak langsung iya sih
		votes <-voteResult{term: rn.Persistence.CurrentTerm, voteGranted: false}
		return
	}

	votes <- voteResult{term: rn.Persistence.CurrentTerm, voteGranted: reply.VoteGranted}
}
