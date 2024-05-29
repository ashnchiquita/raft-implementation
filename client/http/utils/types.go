package utils

type ResponseMessage struct {
	Message string `json:"message"`
}

type KeyVal struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type KeyValResponse struct {
	ResponseMessage
	Data KeyVal `json:"data"`
}

type KeyValsResponse struct {
	ResponseMessage
	Data []KeyVal `json:"data"`
}
