package data_test

import (
	"testing"

	"tubes.sister/raft/node/data"
)

func TestMarshalClusterList(t *testing.T) {
	adddressList := []data.Address{
		{IP: "localhost", Port: 5000},
		{IP: "localhost", Port: 5001},
		{IP: "localhost", Port: 5002},
	}

	marshalledList, err := data.MarshallConfiguration(adddressList)

	expected := "[{\"ip\":\"localhost\",\"port\":5000},{\"ip\":\"localhost\",\"port\":5001},{\"ip\":\"localhost\",\"port\":5002}]"

	if err != nil {
		t.Errorf("Err: %s", err.Error())
	}

	if marshalledList != expected {
		t.Errorf("Marshalled list is incorrect \nexpected: %s \nfound: %s", expected, marshalledList)
	}
}

func TestUnmarshalClusterList(t *testing.T) {
	expectedAdddressList := []data.Address{
		{IP: "localhost", Port: 5000},
		{IP: "localhost", Port: 5001},
		{IP: "localhost", Port: 5002},
	}

	marshalledList := "[{\"ip\":\"localhost\",\"port\":5000},{\"ip\":\"localhost\",\"port\":5001},{\"ip\":\"localhost\",\"port\":5002}]"
	resultingAddressList, err := data.UnmarshallConfiguration(marshalledList)

	if err != nil {
		t.Errorf("Err: %s", err.Error())
	}

	equal := true

	for idx, resultAddress := range resultingAddressList {
		equal = equal && resultAddress.Equals(&expectedAdddressList[idx])
	}

	if !equal {
		t.Errorf("Unmarshalled list is incorrect \nexpected: %v \nfound: %v", expectedAdddressList, resultingAddressList)
	}
}
