package rpc

import (
	"tubes.sister/raft/server"
)

type LogEntry struct {
	key   string
	value string
	term  int
}

type AppendEntriesArgs struct {
	term         int            // leader’s term
	leaderId     server.Address // so follower can redirect clients
	prevLogIndex int            // index of log entry immediately preceding new ones
	prevLogTerm  int            // term of prevLogIndex entry
	entries      []LogEntry     // log entries to store (empty for heartbeat; may send more than one for efficiency)
	leaderCommit int            // leader’s commitIndex
}

type AppendEntriesReply struct {
	term    int  // currentTerm, for leader to update itself
	success bool // true if follower contained entry matching prevLogIndex and prevLogTerm
}

// KEEP IN MIND THAT INDEXES IN PAPER START WITH 1, NOT 0, causing initial index values to be 0
// Current AppendEntries implementation start indexes with 0, causeing initial index values to be -1 (as per comments)

type AppenEntriesTestNode struct {
	// Persistent state on all servers
	// (Updated on stable storage before responding to RPCs)
	currentTerm int            // latest term server has seen (initialized to -1 on first boot, increases monotonically)
	votedFor    server.Address // candidateId that received vote in current term (or null if none)
	log         []LogEntry     // log entries; each entry contains command for state machine, and term when entry was received by leader (first index is 1)

	// Volatile state on all servers
	commitIndex int // index of highest log entry known to be committed (initialized to -1, increases monotonically)
	lastApplied int // index of highest log entry applied to state machine (initialized to -1, increases monotonically)

	// Volatile state on leaders
	// (Reinitialized after election)
	nextIndex  map[server.Address]int // for each server, index of the next log entry to send to that server (initialized to leader last log index + 1)
	matchIndex map[server.Address]int // for each server, index of highest log entry known to be replicated on server (initialized to -1, increases monotonically)

	// Gaada di bagian state di paper
	nodeType     server.NodeType
	stateMachine map[string]string
}

func (node *AppenEntriesTestNode) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) error {
	// Rule 1 : Reply false if term < currentTerm (§5.1)
	if args.term < node.currentTerm {
		reply.success = false
		return nil
	}

	if args.term > node.currentTerm {
		node.currentTerm = args.term
		node.nodeType = server.FOLLOWER
	}

	// Rule 2: Reply false if log doesn’t contain an
	// entry at prevLogIndex whose term matches prevLogTerm
	if len(node.log) <= args.prevLogIndex || node.log[args.prevLogIndex].term != args.prevLogTerm {
		reply.success = false
		return nil
	}

	for i := args.prevLogIndex + 1; i <= args.prevLogIndex+len(args.entries); i++ {
		// Rule 3: If an existing entry conflicts with a new one (same index
		// but different terms), delete the existing entry and all that
		// follow it (§5.3)
		if len(node.log)-1 >= i && node.log[i].term != args.term {
			node.log = node.log[0:i]
		}

		// Rule 4: Append any new entries not already in the log
		if len(node.log) == i {
			node.log = append(node.log, args.entries[i-args.prevLogIndex-1])
		}
	}

	node.commitEntries(args)

	return nil
}

// Rule 5:  If leaderCommit > commitIndex, set commitIndex =
// min(leaderCommit, index of last new entry)
func (node *AppenEntriesTestNode) commitEntries(args *AppendEntriesArgs) error {
	for i := node.commitIndex + 1; i <= args.leaderCommit && i < len(node.log); i++ {
		currentEntry := node.log[i]

		node.stateMachine[currentEntry.key] = currentEntry.value
		node.commitIndex = i
	}

	return nil
}
