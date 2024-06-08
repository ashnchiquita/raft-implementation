package core

type ElectionInterruptMsg int

const (
	ELECTION_TIMEOUT ElectionInterruptMsg = iota + 1
	HIGHER_TERM
)
