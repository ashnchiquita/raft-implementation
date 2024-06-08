package core

import (
	"errors"
	"fmt"
	"strings"

	"tubes.sister/raft/node/data"
)

func (rn *RaftNode) LeaderEnterJointConsensus(marshalledNew string) error {
	oldAddressList := make([]data.Address, len(rn.Volatile.ClusterList))
	for _, clusterData := range rn.Volatile.ClusterList {
		oldAddressList = append(oldAddressList, clusterData.Address)
	}

	marshalledOld, err := data.MarshallConfiguration(oldAddressList)

	if err != nil {
		return err
	}

	newConfig, err := data.UnmarshallConfiguration(marshalledNew)

	if err != nil {
		return err
	}

	newClusterList := data.ClusterListFromAddresses(newConfig, len(rn.Persistence.Log)+1)
	oldNewConfig := data.LogEntry{Term: rn.Persistence.CurrentTerm, Command: "OLDNEWCONF", Value: fmt.Sprintf("%s,%s", marshalledOld, marshalledNew)}
	rn.Persistence.Log = append(rn.Persistence.Log, oldNewConfig)

	rn.Volatile.OldClusterList = rn.Volatile.ClusterList
	rn.Volatile.ClusterList = newClusterList
	rn.Volatile.IsJointConsensus = true
	rn.SetupClusterClients()

	return nil
}

func (rn *RaftNode) FollowerEnterJointConsensus(marshalledOldNew string) error {
	splitMarshalledOldnew := strings.Split(marshalledOldNew, ",")

	if len(splitMarshalledOldnew) != 2 {
		return errors.New("follower received config entry with invalid format")
	}

	marshalledOld := splitMarshalledOldnew[0]
	marshalledNew := splitMarshalledOldnew[1]

	oldAddressList, err := data.UnmarshallConfiguration(marshalledOld)

	if err != nil {
		return err
	}

	newAddressList, err := data.UnmarshallConfiguration(marshalledNew)

	if err != nil {
		return err
	}

	oldClusterList := data.ClusterListFromAddresses(oldAddressList, len(rn.Persistence.Log)+1)
	newClusterList := data.ClusterListFromAddresses(newAddressList, len(rn.Persistence.Log)+1)

	rn.Volatile.OldClusterList = oldClusterList
	rn.Volatile.ClusterList = newClusterList
	rn.Volatile.IsJointConsensus = true
	rn.SetupClusterClients()

	return nil
}

func (rn *RaftNode) ApplyNewClusterList(marshalledNew string) error {
	newAddressList, err := data.UnmarshallConfiguration(marshalledNew)

	if err != nil {
		return err
	}

	newClusterList := data.ClusterListFromAddresses(newAddressList, len(rn.Persistence.Log)+1)

	rn.Volatile.OldClusterList = nil
	rn.Volatile.ClusterList = newClusterList
	rn.Volatile.IsJointConsensus = false
	rn.SetupClusterClients()

	return nil
}
