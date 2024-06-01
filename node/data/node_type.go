package data

type NodeType int

const (
	LEADER NodeType = iota + 1
	CANDIDATE
	FOLLOWER
)
