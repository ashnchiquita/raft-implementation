package data

import (
	"reflect"
	"strconv"
)

type Address struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

// CONSTRUCTOR
func NewAddress(ip string, port int) *Address {
	return &Address{
		IP:   ip,
		Port: port,
	}
}

func NewZeroAddress() *Address {
	return &Address{
		IP:   "",
		Port: 0,
	}
}

func (a *Address) String() string {
	return a.IP + ":" + strconv.Itoa(a.Port)
}

func (a *Address) Equals(other *Address) bool {
	return reflect.DeepEqual(a, other)
}

func (a *Address) IsZeroAddress() bool {
	return a.IP == "" && a.Port == 0
}

func UnionAddressList(a []Address, b []Address) []Address {
	addressMap := make(map[Address]bool)

	for _, address := range a {
		addressMap[address] = true
	}

	for _, address := range b {
		addressMap[address] = true
	}

	keys := make([]Address, len(addressMap))

	i := 0
	for k := range addressMap {
		keys[i] = k
		i++
	}

	return keys
}
