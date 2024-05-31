package data

import (
	"reflect"
)

type Address struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

// CONSTRUCTOR
func NewAddress(ip, port string) *Address {
	return &Address{
		IP: ip,
		Port: port,
	}
}

func (a *Address) String() string {
	return a.IP + ":" + a.Port
}

func (a *Address) Equals(other *Address) bool {
	return reflect.DeepEqual(a, other)
}
