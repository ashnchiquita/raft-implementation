package data

type Volatile struct {
	LeaderAddress Address
	ClusterList   []ClusterData
	CommitIndex   int
	LastApplied   int
	Type          NodeType
	VotesReceived map[Address]bool
}

// CONSTRUCTOR
func NewVolatile() *Volatile {
	return &Volatile{
		CommitIndex: -1,
		LastApplied: -1,
		Type:        FOLLOWER,
		VotesReceived: make(map[Address]bool),
	}
}

func (v *Volatile) ResetVotes() {
	v.VotesReceived = make(map[Address]bool)
}

func (v *Volatile) GetVotesCount() int {
	return len(v.VotesReceived)
}

func (v *Volatile) AddVote(address Address) {
	v.VotesReceived[address] = true
}
