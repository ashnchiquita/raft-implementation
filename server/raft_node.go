package server

type RaftNode struct {
	address           Address
	nodeType          NodeType
	log               []string
	electionTerm      int
	clusterAddrList   []Address
	clusterLeaderAddr Address
}
