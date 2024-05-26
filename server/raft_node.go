package server

import "log"

type RaftNode struct {
	address           Address
	nodeType          NodeType
	log               []string
	electionTerm      int
	clusterAddrList   []Address
	clusterLeaderAddr Address
}

func NewRaftNode(address Address) *RaftNode {
	return &RaftNode{
		address:           address,
		nodeType:          0,
		log:               []string{},
		electionTerm:      0,
		clusterAddrList:   []Address{},
		clusterLeaderAddr: address, // saat buat node baru leadernya diri sendiri(?)
	}
}

func (rn *RaftNode) AddNode(newNodeAddress Address) {
	rn.clusterAddrList = append(rn.clusterAddrList, newNodeAddress)
	// DEBUGING
	log.Println("Updated clusterAddrList:", rn.clusterAddrList)
}
