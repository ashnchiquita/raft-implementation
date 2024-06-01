package data

type Persistence struct {
	CurrentTerm int        `json:"currentTerm"`
	VotedFor    Address    `json:"votedFor"`
	Log         []LogEntry `json:"log"`
}

// CONSTRUCTOR
func NewPersistence() *Persistence {
	// initialized votedFor to nil
	return &Persistence{
		CurrentTerm: 0,
		Log:         []LogEntry{},
	}
}
