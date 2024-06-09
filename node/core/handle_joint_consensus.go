package core

import (
	"encoding/json"
	"log"
	"os"

	"tubes.sister/raft/node/data"
)

func (rn *RaftNode) LeaderEnterJointConsensus(marshalledNew string) error {
	log.Printf("Memcange >> Leader is entering joint consensus; newconfig: %s", marshalledNew)
	newConfig, err := data.UnmarshallConfiguration(marshalledNew)

	if err != nil {
		return err
	}

	rn.Volatile.OldConfig = data.ClusterListToAddressList(rn.Volatile.ClusterList)
	rn.Volatile.NewConfig = newConfig

	data.DisconnectClusterList(rn.Volatile.ClusterList)
	rn.Volatile.ClusterList = data.ClusterListFromAddresses(
		data.UnionAddressList(
			rn.Volatile.NewConfig,
			rn.Volatile.OldConfig),
		len(rn.Persistence.Log))

	payload := data.OldNewConfigPayload{New: rn.Volatile.NewConfig, Old: rn.Volatile.OldConfig}
	marshalledOldNew, err := payload.Marshall()

	if err != nil {
		return err
	}

	oldNewConfig := data.LogEntry{Term: rn.Persistence.CurrentTerm, Command: "OLDNEWCONF", Value: marshalledOldNew}
	rn.Persistence.Log = append(rn.Persistence.Log, oldNewConfig)

	rn.Volatile.IsJointConsensus = true

	return nil
}

func (rn *RaftNode) FollowerEnterJointConsensus(marshalledOldNew string) error {
	log.Printf("Memcange >> Follower is entering joint consensus; oldnewconfig: %s", marshalledOldNew)

	oldNewConfig, err := data.NewOldConfigPayloadFromJson(marshalledOldNew)

	log.Println(oldNewConfig)

	if err != nil {
		return err
	}

	rn.Volatile.OldConfig = oldNewConfig.Old
	rn.Volatile.NewConfig = oldNewConfig.New

	data.DisconnectClusterList(rn.Volatile.ClusterList)
	rn.Volatile.ClusterList = data.ClusterListFromAddresses(
		data.UnionAddressList(
			rn.Volatile.OldConfig,
			rn.Volatile.NewConfig),
		len(rn.Persistence.Log))

	rn.Volatile.IsJointConsensus = true

	return nil
}

func (rn *RaftNode) ApplyNewClusterList(marshalledNew string) error {
	log.Printf("Memcange >> Applying new config: %s", marshalledNew)
	newAddressList, err := data.UnmarshallConfiguration(marshalledNew)

	if err != nil {
		return err
	}

	newClusterList := data.ClusterListFromAddresses(newAddressList, len(rn.Persistence.Log))

	rn.Volatile.OldConfig = nil
	rn.Volatile.NewConfig = nil

	data.DisconnectClusterList(rn.Volatile.ClusterList)
	rn.Volatile.ClusterList = newClusterList
	rn.Volatile.IsJointConsensus = false

	return nil
}

func (rn *RaftNode) ReadClusterConfigFromFile(path string) error {
	var addressList []data.Address
	file, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &addressList)

	data.DisconnectClusterList(rn.Volatile.ClusterList)
	rn.Volatile.ClusterList = data.ClusterListFromAddresses(addressList, len(rn.Persistence.Log))

	return err
}
