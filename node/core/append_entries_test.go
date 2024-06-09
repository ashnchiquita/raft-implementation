package core_test

import (
	"context"
	"log"
	"testing"

	pb "tubes.sister/raft/gRPC/node/core"
	"tubes.sister/raft/node/core"
	"tubes.sister/raft/node/data"
)

func CreateServerTest(isLeader bool, port int, initialize bool, logs []data.LogEntry) *core.RaftNode {
	addr := data.NewAddress("localhost", port)

	raftNode := core.NewRaftNode(*addr)
	if isLeader {
		raftNode.Volatile.Type = data.LEADER
	}
	raftNode.InitializeAsLeader()
	raftNode.Volatile.ClusterList = []data.ClusterData{
		*data.NewClusterData(*data.NewAddress("localhost", 5000), 0, -1),
		*data.NewClusterData(*data.NewAddress("localhost", 5001), 0, -1),
		*data.NewClusterData(*data.NewAddress("localhost", 5002), 0, -1),
	}

	raftNode.Volatile.LeaderAddress = *data.NewAddress("localhost", 5000)
	raftNode.Persistence.Log = logs
	if initialize {
		raftNode.InitializeServer()
	}

	return raftNode
}

func TestAppendEntriesRejectLowerTerm(t *testing.T) {
	logs := []data.LogEntry{}
	rn := CreateServerTest(false, 5001, false, logs)

	rn.Persistence.CurrentTerm = 2

	args := pb.AppendEntriesArgs{
		Term: int32(1),
		LeaderAddress: &pb.AppendEntriesArgs_LeaderAddress{
			Ip:   "localhost",
			Port: 5000,
		},
		PrevLogIndex: 0,
		PrevLogTerm:  1,
		Entries:      []*pb.AppendEntriesArgs_LogEntry{},
		LeaderCommit: 0,
	}

	reply, err := rn.AppendEntries(context.Background(), &args)

	if err != nil {
		t.Errorf("Err: %s", err.Error())
	}

	if reply.Success {
		t.Error("Append entries does not reject leader with lower term")
	}
}

func TestAppendEntriesRevertToFollowerOnHigherTerm(t *testing.T) {
	logs := []data.LogEntry{
		*data.NewLogEntry(1, "SET", data.WithValue("1")),
	}
	rn := CreateServerTest(false, 5001, false, logs)

	rn.Volatile.Type = data.LEADER
	rn.Persistence.CurrentTerm = 1

	args := pb.AppendEntriesArgs{
		Term: int32(2),
		LeaderAddress: &pb.AppendEntriesArgs_LeaderAddress{
			Ip:   "localhost",
			Port: 5000,
		},
		PrevLogIndex: 0,
		PrevLogTerm:  3,
		Entries:      []*pb.AppendEntriesArgs_LogEntry{},
		LeaderCommit: 0,
	}

	_, err := rn.AppendEntries(context.Background(), &args)

	if err != nil {
		t.Errorf("Err: %s", err.Error())
	}

	if rn.Volatile.Type != data.FOLLOWER {
		t.Error("Append entries doesn't turn leader into follower when a leader received append entries with higher term")
	}
}

func TestAppendEntriesRejectNonMatchingPrevLogIndex(t *testing.T) {
	logs := []data.LogEntry{
		*data.NewLogEntry(1, "SET", data.WithValue("1")),
	}
	rn := CreateServerTest(false, 5001, false, logs)

	rn.Persistence.CurrentTerm = 1

	args := pb.AppendEntriesArgs{
		Term: int32(1),
		LeaderAddress: &pb.AppendEntriesArgs_LeaderAddress{
			Ip:   "localhost",
			Port: 5000,
		},
		PrevLogIndex: 0,
		PrevLogTerm:  3,
		Entries:      []*pb.AppendEntriesArgs_LogEntry{},
		LeaderCommit: 0,
	}

	reply, err := rn.AppendEntries(context.Background(), &args)

	if err != nil {
		t.Errorf("Err: %s", err.Error())
	}

	if reply.Success {
		t.Error("Append entries doesn't fail non matching prev log index")
	}
}

func TestAppendEntriesClearConflictingEntries(t *testing.T) {
	logs := []data.LogEntry{
		*data.NewLogEntry(1, "SET", data.WithValue("1")),
		*data.NewLogEntry(1, "SET", data.WithValue("1")),
		*data.NewLogEntry(2, "SET", data.WithValue("1")),
		*data.NewLogEntry(2, "SET", data.WithValue("1")),
		*data.NewLogEntry(2, "SET", data.WithValue("1")),
	}
	rn := CreateServerTest(false, 5001, false, logs)

	rn.Persistence.CurrentTerm = 3

	args := pb.AppendEntriesArgs{
		Term: int32(3),
		LeaderAddress: &pb.AppendEntriesArgs_LeaderAddress{
			Ip:   "localhost",
			Port: 5000,
		},
		PrevLogIndex: 2,
		PrevLogTerm:  2,
		Entries:      []*pb.AppendEntriesArgs_LogEntry{{Term: 1, Command: "SET", Value: "1"}},
		LeaderCommit: -1,
	}

	reply, err := rn.AppendEntries(context.Background(), &args)
	log.Print(reply.Success)

	if err != nil {
		t.Errorf("Err: %s", err.Error())
	}

	if len(rn.Persistence.Log) != 4 {
		t.Errorf("Append entries doesn't clear conflicting entries, log length: %d", len(rn.Persistence.Log))
	}
}
func TestAppendEntriesAppendNewEntries(t *testing.T) {
	logs := []data.LogEntry{
		*data.NewLogEntry(1, "SET", data.WithValue("1")),
		*data.NewLogEntry(1, "SET", data.WithValue("1")),
		*data.NewLogEntry(2, "SET", data.WithValue("1")),
		*data.NewLogEntry(2, "SET", data.WithValue("1")),
		*data.NewLogEntry(2, "SET", data.WithValue("1")),
	}
	rn := CreateServerTest(false, 5001, false, logs)

	rn.Persistence.CurrentTerm = 3

	args := pb.AppendEntriesArgs{
		Term: int32(3),
		LeaderAddress: &pb.AppendEntriesArgs_LeaderAddress{
			Ip:   "localhost",
			Port: 5000,
		},
		PrevLogIndex: 4,
		PrevLogTerm:  2,
		Entries: []*pb.AppendEntriesArgs_LogEntry{
			{Term: 3, Command: "SET", Value: "1"},
			{Term: 3, Command: "SET", Value: "1"},
			{Term: 3, Command: "SET", Value: "1"},
		},
		LeaderCommit: -1,
	}

	reply, err := rn.AppendEntries(context.Background(), &args)
	log.Print(reply.Success)

	if err != nil {
		t.Errorf("Err: %s", err.Error())
	}

	if len(rn.Persistence.Log) != 8 {
		t.Errorf("Append entries fails to add new logs, log length: %d", len(rn.Persistence.Log))
	}
}

func TestAppendEntriesCommitEntries(t *testing.T) {
	logs := []data.LogEntry{
		*data.NewLogEntry(1, "SET", data.WithValue("1,1")),
		*data.NewLogEntry(1, "SET", data.WithValue("2,2")),
		*data.NewLogEntry(2, "APPEND", data.WithValue("1, 3")),
		*data.NewLogEntry(2, "DEL", data.WithValue("2")),
	}
	rn := CreateServerTest(false, 5001, false, logs)

	rn.Persistence.CurrentTerm = 3

	args := pb.AppendEntriesArgs{
		Term: int32(3),
		LeaderAddress: &pb.AppendEntriesArgs_LeaderAddress{
			Ip:   "localhost",
			Port: 5000,
		},
		PrevLogIndex: 4,
		PrevLogTerm:  2,
		Entries:      []*pb.AppendEntriesArgs_LogEntry{},
		LeaderCommit: 4,
	}

	reply, err := rn.AppendEntries(context.Background(), &args)
	log.Print(reply.Success)

	if err != nil {
		t.Errorf("Err: %s", err.Error())
	}

	data := rn.Application.GetAll()

	if data["1"] != "1 3" && len(data) == 1 {
		t.Errorf("Append entries failed commit entries, state machine: %v", data)
	}
}
