package data

type Volatile struct {
	CommitIndex   int
	LastApplied   int
	LeaderAddress Address
	ClusterList   []Address
}

// CONSTRUCTOR
func NewVolatile() *Volatile {
	return &Volatile{
		CommitIndex: 0,
		LastApplied: 0,
	}
}
