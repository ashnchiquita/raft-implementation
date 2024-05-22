package server

type NodeType int

const (
	LEADER NodeType = iota + 1
	CANDIDATE
	FOLLOWER
)
