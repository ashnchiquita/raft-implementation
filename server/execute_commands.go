package server

type ExecuteCommand string

const (
	PING   ExecuteCommand = "Ping"
	GET    ExecuteCommand = "Get"
	SET    ExecuteCommand = "Set"
	STRLEN ExecuteCommand = "Strlen"
	DEL    ExecuteCommand = "Del"
	APPEND ExecuteCommand = "Append"

	// NICE TO HAVE METHODS
	GETALL ExecuteCommand = "GetAll"
	DELALL ExecuteCommand = "DelAll"
)
