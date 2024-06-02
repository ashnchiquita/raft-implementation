package data

type Volatile struct {
	LeaderAddress Address
	ClusterList   []ClusterData
	CommitIndex   int
	LastApplied   int
	Type          NodeType
}

// CONSTRUCTOR
func NewVolatile() *Volatile {
	return &Volatile{
		CommitIndex: -1,
		LastApplied: -1,
		Type:        FOLLOWER,
	}
}
