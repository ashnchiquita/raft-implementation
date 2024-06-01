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

func (a *Address) String() string {
	return a.IP + ":" + strconv.Itoa(a.Port)
}

func (a *Address) Equals(other *Address) bool {
	return reflect.DeepEqual(a, other)
}
