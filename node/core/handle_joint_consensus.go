package core

import (
	"encoding/json"
	"log"
	"os"

	"tubes.sister/raft/node/data"
)

func (rn *RaftNode) LeaderEnterJointConsensus(marshalledNew string) error {
	rn.logf(Purple+"Memchange >> "+Reset+"Leader is entering joint consensus; newconfig: %s", marshalledNew)
	newConfig, err := data.UnmarshallConfiguration(marshalledNew)

	if err != nil {
		return err
	}

	rn.Volatile.OldConfig = data.ClusterListToAddressList(rn.Volatile.ClusterList)
	rn.Volatile.NewConfig = newConfig

	data.DisconnectClusterList(rn.Volatile.ClusterList)

	clusterDataMap := make(map[data.Address]data.ClusterData)

	for idx, clusterData := range rn.Volatile.ClusterList {
		clusterDataMap[clusterData.Address] = rn.Volatile.ClusterList[idx]
	}

	rn.Volatile.ClusterList = data.ClusterListFromAddresses(
		data.UnionAddressList(
			rn.Volatile.NewConfig,
			rn.Volatile.OldConfig),
		len(rn.Persistence.Log))

	for idx, node := range rn.Volatile.ClusterList {
		if _, ok := clusterDataMap[node.Address]; !ok {
			continue
		}

		rn.Volatile.ClusterList[idx].MatchIndex = clusterDataMap[node.Address].MatchIndex
		rn.Volatile.ClusterList[idx].NextIndex = clusterDataMap[node.Address].NextIndex
	}

	payload := data.OldNewConfigPayload{New: rn.Volatile.NewConfig, Old: rn.Volatile.OldConfig}
	marshalledOldNew, err := payload.Marshall()

	if err != nil {
		return err
	}

	oldNewConfig := data.LogEntry{Term: rn.Persistence.CurrentTerm, Command: "OLDNEWCONF", Value: marshalledOldNew}
	rn.Persistence.Log = append(rn.Persistence.Log, oldNewConfig)
	rn.Persistence.Serialize()

	rn.Volatile.IsJointConsensus = true

	return nil
}

func (rn *RaftNode) FollowerEnterJointConsensus(marshalledOldNew string) error {
	rn.logf(Purple+"Memchange >> "+Reset+"Follower is entering joint consensus; oldnewconfig: %s", marshalledOldNew)

	oldNewConfig, err := data.NewOldConfigPayloadFromJson(marshalledOldNew)

	log.Println(oldNewConfig)

	if err != nil {
		return err
	}

	rn.Volatile.OldConfig = oldNewConfig.Old
	rn.Volatile.NewConfig = oldNewConfig.New

	clusterDataMap := make(map[data.Address]data.ClusterData)

	for idx, clusterData := range rn.Volatile.ClusterList {
		clusterDataMap[clusterData.Address] = rn.Volatile.ClusterList[idx]
	}

	data.DisconnectClusterList(rn.Volatile.ClusterList)
	rn.Volatile.ClusterList = data.ClusterListFromAddresses(
		data.UnionAddressList(
			rn.Volatile.OldConfig,
			rn.Volatile.NewConfig),
		len(rn.Persistence.Log))

	for idx, node := range rn.Volatile.ClusterList {
		if _, ok := clusterDataMap[node.Address]; !ok {
			continue
		}

		rn.Volatile.ClusterList[idx].MatchIndex = clusterDataMap[node.Address].MatchIndex
		rn.Volatile.ClusterList[idx].NextIndex = clusterDataMap[node.Address].NextIndex
	}

	rn.Volatile.IsJointConsensus = true

	return nil
}

func (rn *RaftNode) ApplyNewClusterList(marshalledNew string) error {
	rn.logf(Purple+"Memchange >> "+Reset+"Applying new config: %s", marshalledNew)
	newAddressList, err := data.UnmarshallConfiguration(marshalledNew)

	if err != nil {
		return err
	}

	newClusterList := data.ClusterListFromAddresses(newAddressList, len(rn.Persistence.Log))

	clusterDataMap := make(map[data.Address]data.ClusterData)

	rn.Volatile.OldConfig = nil
	rn.Volatile.NewConfig = nil

	for idx, clusterData := range rn.Volatile.ClusterList {
		clusterDataMap[clusterData.Address] = rn.Volatile.ClusterList[idx]
	}

	data.DisconnectClusterList(rn.Volatile.ClusterList)
	rn.Volatile.ClusterList = newClusterList
	rn.Volatile.IsJointConsensus = false

	for idx, node := range rn.Volatile.ClusterList {
		if _, ok := clusterDataMap[node.Address]; !ok {
			continue
		}

		rn.Volatile.ClusterList[idx].MatchIndex = clusterDataMap[node.Address].MatchIndex
		rn.Volatile.ClusterList[idx].NextIndex = clusterDataMap[node.Address].NextIndex
	}

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
