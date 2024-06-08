package data

import "encoding/json"

type Volatile struct {
	LeaderAddress    Address
	ClusterList      []ClusterData
	CommitIndex      int
	LastApplied      int
	Type             NodeType
	IsJointConsensus bool
	OldClusterList   []ClusterData
	VotesReceived    map[Address]bool
}

// CONSTRUCTOR
func NewVolatile() *Volatile {
	return &Volatile{
		CommitIndex:   -1,
		LastApplied:   -1,
		Type:          FOLLOWER,
		VotesReceived: make(map[Address]bool),
	}
}

func MarshallConfiguration(addressList []Address) (string, error) {
	b, err := json.Marshal(addressList)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func UnmarshallConfiguration(marshalledList string) ([]Address, error) {
	var data []Address
	err := json.Unmarshal([]byte(marshalledList), &data)

	return data, err
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

func (v *Volatile) GetVoters() []Address {
	voters := make([]Address, 0)
	for k := range v.VotesReceived {
		voters = append(voters, k)
	}
	return voters
}
