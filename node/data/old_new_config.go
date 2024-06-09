package data

import "encoding/json"

type OldNewConfigPayload struct {
	Old []Address `json:"old"`
	New []Address `json:"new"`
}

func (p *OldNewConfigPayload) Marshall() (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func NewOldConfigPayloadFromJson(marshalledConfig string) (*OldNewConfigPayload, error) {
	var data OldNewConfigPayload
	err := json.Unmarshal([]byte(marshalledConfig), &data)

	return &data, err
}
