package core

import (
	"context"
	"testing"

	gRPC "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/data"
)

// func TestFooer(t *testing.T) {
//     result := Fooer(3)
//     if result != "Foo" {
//     t.Errorf("Result was incorrect, got: %s, want: %s.", result, "Foo")
//     }
// }

func TestRequestVote(t *testing.T) {

	rn := &RaftNode{
		Persistence: data.Persistence{
			Address:     data.Address{},
			CurrentTerm: 1,
			VotedFor:    data.Address{},
			Log:         []data.LogEntry{{Term: 1, Command: "SET"}},
		},
		Volatile: data.Volatile{
			LeaderAddress: *data.NewAddress("localhost", 5000),
			ClusterList: []data.ClusterData{
				{Address: *data.NewAddress("localhost", 5000), MatchIndex: -1, NextIndex: 0},
				{Address: *data.NewAddress("localhost", 5001), MatchIndex: -1, NextIndex: 0},
				{Address: *data.NewAddress("localhost", 5002), MatchIndex: -1, NextIndex: 0},
			},
			CommitIndex: 0,
			LastApplied: 0,
			Type:        data.LEADER,
		},
	}

	tests := []struct {
		name            string
		args            *gRPC.RequestVoteArgs
		wantTerm        int32
		wantVoteGranted bool
	}{
		{
			name: "term < currentTerm",
			args: &gRPC.RequestVoteArgs{
				Term:             0,
				CandidateAddress: &gRPC.RequestVoteArgs_CandidateAddress{},
				LastLogIndex:     1,
				LastLogTerm:      1,
			},
			wantTerm:        1,
			wantVoteGranted: false,
		},
		{
			name: "term > currentTerm, updated log",
			args: &gRPC.RequestVoteArgs{
				Term:             2,
				CandidateAddress: &gRPC.RequestVoteArgs_CandidateAddress{},
				LastLogIndex:     1,
				LastLogTerm:      1,
			},
			wantTerm:        2,
			wantVoteGranted: true,
		},
		{
			name: "term > currentTerm, not updated log",
			args: &gRPC.RequestVoteArgs{
				Term:             2,
				CandidateAddress: &gRPC.RequestVoteArgs_CandidateAddress{},
				LastLogIndex:     0,
				LastLogTerm:      0,
			},
			wantTerm:        2,
			wantVoteGranted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rn.RequestVote(context.Background(), tt.args)
			if err != nil {
				t.Fatalf("RequestVote() error = %v", err)
			}
			if got.Term != tt.wantTerm {
				t.Errorf("RequestVote() got term = %v, want %v", got.Term, tt.wantTerm)
			}
			if got.VoteGranted != tt.wantVoteGranted {
				t.Errorf("RequestVote() got vote = %v, want %v", got.VoteGranted, tt.wantVoteGranted)
			}
		})
	}
}
