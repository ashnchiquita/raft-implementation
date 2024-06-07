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

func MarshallClusterList(clusterList []ClusterData) (string, error) {
	b, err := json.Marshal(clusterList)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func UnmarshallClusterList(marshalledList string) ([]ClusterData, error) {
	var data []ClusterData
	err := json.Unmarshal([]byte(marshalledList), &data)

	return data, err
}
