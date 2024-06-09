package data

import (
	"encoding/json"

	"tubes.sister/raft/node/utils/serializer"
)

type Persistence struct {
	Address 		Address    `json:"address"`
	CurrentTerm int        `json:"currentTerm"`
	VotedFor    Address    `json:"votedFor"`
	Log         []LogEntry `json:"log"`
}

// SERIALIZATION
// TODO: change to gob if everything works
var ser = serializer.NewJSONSerializer[Persistence]()

// CONSTRUCTOR: will initialize votedFor to Zero Value
func NewPersistence(address Address) *Persistence {
	return &Persistence{
		Address:     address,
		CurrentTerm: 0,
		Log:         []LogEntry{},
	}
}

func (p *Persistence) Serialize() error {
	return ser.Serialize(*p, p.Address.String())
}

func Deserialize(address Address) (*Persistence, error) {
	data, err := ser.Deserialize(address.String())
	return &data, err
}

// Will try to load from file, if not found/error, create a new one (not persisted yet)
func InitPersistence(address Address) *Persistence {
	data, err := Deserialize(address)
	if err != nil {
		data = NewPersistence(address)
	}

	return data
}

func (p *Persistence) GetPrettyLog() (string, error) {
	bytes, err := json.MarshalIndent(p.Log, "", "  ")

	return string(bytes), err
}
