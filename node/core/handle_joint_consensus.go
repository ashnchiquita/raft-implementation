package core

import (
	"errors"
	"fmt"
	"strings"

	"tubes.sister/raft/node/data"
)

func (rn *RaftNode) LeaderEnterJointConsensus(marshalledNew string) error {
	marshalledOld, err := data.MarshallClusterList(rn.Volatile.ClusterList)

	if err != nil {
		return err
	}

	newClusterList, err := data.UnmarshallClusterList(marshalledNew)

	if err != nil {
		return err
	}

	oldNewConfig := data.LogEntry{Term: rn.Persistence.CurrentTerm, Command: "OLDNEWCONF", Value: fmt.Sprintf("%s,%s", marshalledOld, marshalledNew)}
	rn.Persistence.Log = append(rn.Persistence.Log, oldNewConfig)

	rn.Volatile.OldClusterList = rn.Volatile.ClusterList
	rn.Volatile.ClusterList = newClusterList
	rn.Volatile.IsJointConsensus = true

	return nil
}

func (rn *RaftNode) FollowerEnterJointConsensus(marshalledOldNew string) error {
	splitMarshalledOldnew := strings.Split(marshalledOldNew, ",")

	if len(splitMarshalledOldnew) != 2 {
		return errors.New("follower received config entry with invalid format")
	}

	marshalledOld := splitMarshalledOldnew[0]
	marshalledNew := splitMarshalledOldnew[1]

	oldClusterList, err := data.UnmarshallClusterList(marshalledOld)

	if err != nil {
		return err
	}

	newClusterList, err := data.UnmarshallClusterList(marshalledNew)

	if err != nil {
		return err
	}

	rn.Volatile.OldClusterList = oldClusterList
	rn.Volatile.ClusterList = newClusterList
	rn.Volatile.IsJointConsensus = true

	return nil
}

func (rn *RaftNode) ApplyNewClusterList(marshalledNew string) error {
	newClusterList, err := data.UnmarshallClusterList(marshalledNew)

	if err != nil {
		return err
	}

	rn.Volatile.OldClusterList = nil
	rn.Volatile.ClusterList = newClusterList
	rn.Volatile.IsJointConsensus = false

	return nil
}
