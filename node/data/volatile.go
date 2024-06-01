package data

type Volatile struct {
	LeaderAddress Address
	ClusterList   []Address
	CommitIndex   int
	LastApplied   int
	Type          NodeType
}

// CONSTRUCTOR
func NewVolatile() *Volatile {
	return &Volatile{
		CommitIndex: 0,
		LastApplied: 0,
		Type:        FOLLOWER,
	}
}
