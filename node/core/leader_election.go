package core

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/data"
)

type voteResult struct {
	term        int
	voteGranted bool
	address     data.Address
	sameCluster bool
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
	restartElection := make(chan bool)

	go rn.election(restartElection)
	for <-restartElection {
		rn.logf(Yellow + "startElection() >>" + Reset + "Restarting election because of timeout")
		go rn.election(restartElection)
	}
}

// Will convert the current node to a candidate (+ reset its timeout) and start the election process
func (rn *RaftNode) election(restartElection chan bool) {
	rn.announcef("================= STARTING ELECTION ====================")
	logrn := Yellow + "election() >> " + Reset

	var result string
	for _, cluster := range rn.Volatile.ClusterList {
		result += fmt.Sprintf("\n%s", &cluster.Address)
	}
	rn.logf(logrn+"cluster: "+Cyan+"%v"+Reset, result)
	votes := make(chan voteResult)

	rn.setAsCandidate()

	rn.Persistence.CurrentTerm++
	rn.logf(logrn+"Current term: %d", rn.Persistence.CurrentTerm)
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

	// Persist before requesting votes
	rn.Persistence.Serialize()

	for _, node := range rn.Volatile.ClusterList {
		if node.Address == rn.Address {
			continue
		}

		go rn.requestVote(&voteReq, &node, votes)
	}

	requestedToOtherCluster := false

	for {
		// accept reply until interrupted
		select {
		case val := <-rn.electionInterrupt:
			switch val {
			case ELECTION_TIMEOUT:
				rn.cleanupCandidateState()
				rn.logf(logrn + "Election interrupted because of timeout")
				if requestedToOtherCluster {
					data.DisconnectClusterList(rn.Volatile.ClusterList)
					log.Fatalln("election() >> Candidate isn't included in the cluster anymore")
				} else {
					restartElection <- true
				}
			case HIGHER_TERM:
				rn.logf(logrn + "Election interrupted because of higher term")
				restartElection <- false
			}

			return
		case vote := <-votes:
			rn.logf(logrn + "Vote received from channel")
			requestedToOtherCluster = requestedToOtherCluster || !vote.sameCluster

			if !vote.sameCluster {
				rn.logf(logrn + "Requested vote to other cluster")
			}

			if rn.Volatile.Type == data.CANDIDATE && vote.term == rn.Persistence.CurrentTerm && vote.voteGranted {
				rn.logf(logrn+"Vote granted from ", vote.address.String())
				rn.Volatile.AddVote(vote.address)
				rn.logf(logrn+"Votes count: %d", rn.Volatile.GetVotesCount(), " Cluster count: ", len(rn.Volatile.ClusterList))
				rn.logf(logrn+"Voters: ", rn.Volatile.GetVoters())

				if (rn.Volatile.IsJointConsensus &&
					data.MajorityVotedInCluster(rn.Volatile.ClusterList, rn.Volatile.OldConfig, rn.Volatile.VotesReceived, rn.Address) &&
					data.MajorityVotedInCluster(rn.Volatile.ClusterList, rn.Volatile.NewConfig, rn.Volatile.VotesReceived, rn.Address)) ||
					(!rn.Volatile.IsJointConsensus &&
						data.MajorityVotedInCluster(rn.Volatile.ClusterList, data.ClusterListToAddressList(rn.Volatile.ClusterList), rn.Volatile.VotesReceived, rn.Address)) {
					rn.logf(logrn + "Majority reached")
					rn.cleanupCandidateState()
					rn.InitializeAsLeader() // will also reset timeout
					restartElection <- false
					return
				}
			} else if vote.term > rn.Persistence.CurrentTerm {
				rn.logf(logrn+"Higher term received (vote term: ", vote.term, ", curr term: ", rn.Persistence.CurrentTerm, ")")
				rn.logf(logrn + "Setting as follower")
				rn.cleanupCandidateState()
				rn.Persistence.CurrentTerm = vote.term
				rn.setAsFollower() // will also reset timeout
				restartElection <- false
				return
			}
		}
	}
}

func (rn *RaftNode) requestVote(voteReq *voteRequest, node *data.ClusterData, votes chan voteResult) {
	logrn := Green + "RequestVote() >> " + Reset
	rn.logf(logrn+"Requesting vote to to %v", node.Address.Port)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", node.Address.IP, node.Address.Port), opts...)
	if err != nil {
		log.Fatalf("requestVote() >> Failed to dial server: %v", err)
	}
	defer conn.Close()

	client := gRPC.NewRaftNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), RPC_TIMEOUT)

	defer cancel()

	reply, err := client.RequestVote(ctx, &gRPC.RequestVoteArgs{
		Term:             int32(rn.Persistence.CurrentTerm),
		CandidateAddress: &gRPC.RequestVoteArgs_CandidateAddress{Ip: voteReq.candidateAddress.IP, Port: int32(voteReq.candidateAddress.Port)},
		LastLogIndex:     int32(voteReq.lastIndex),
		LastLogTerm:      int32(voteReq.lastTerm),
	})

	if err != nil {
		rn.logf(logrn+"Error requesting vote to %v: %v", node.Address, err)
		return
	}

	res := voteResult{term: int(reply.Term), voteGranted: reply.VoteGranted, address: node.Address, sameCluster: reply.SameCluster}
	rn.logf(logrn+"Vote result: %+v", res)
	votes <- res
}
