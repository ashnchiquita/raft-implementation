package data

type LeaderVolatile struct {
	nextIndex 	[]int
	matchIndex 	[]int
	Volatile
}

// CONSTRUCTOR
func NewLeaderVolatile(commitIndex, lastApplied int) *LeaderVolatile {
	return &LeaderVolatile{
		nextIndex: []int{},
		matchIndex: []int{},
		Volatile: Volatile{
			CommitIndex: commitIndex,
			LastApplied: lastApplied,
		},
	}
}
