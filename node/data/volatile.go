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
}

// CONSTRUCTOR
func NewVolatile() *Volatile {
	return &Volatile{
		CommitIndex: -1,
		LastApplied: -1,
		Type:        FOLLOWER,
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
