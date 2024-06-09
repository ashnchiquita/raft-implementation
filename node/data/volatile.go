package data

import (
	"encoding/json"
)

type Volatile struct {
	LeaderAddress    Address
	ClusterList      []ClusterData
	CommitIndex      int
	LastApplied      int
	Type             NodeType
	IsJointConsensus bool
	OldConfig        []Address
	NewConfig        []Address
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

func MajorityVotedInCluster(clusterList []ClusterData, population []Address, votesReceived map[Address]bool, candidateAddress Address) bool {
	count := 0
	candidateInCluster := 0

	addressMap := make(map[Address]bool)

	for _, address := range population {
		addressMap[address] = true
	}

	for _, clusterData := range clusterList {
		if _, ok := addressMap[clusterData.Address]; !ok {
			continue
		}

		if clusterData.Address.Equals(&candidateAddress) {
			candidateInCluster++
			continue
		}

		if votesReceived[clusterData.Address] {
			count++
		}
	}
	threshold := len(clusterList)/2 + 1
	return count+candidateInCluster >= threshold
}

func ClusterListToAddressList(clusterList []ClusterData) []Address {
	addresses := make([]Address, len(clusterList))

	for idx, clusterData := range clusterList {
		addresses[idx] = clusterData.Address
	}

	return addresses
}
