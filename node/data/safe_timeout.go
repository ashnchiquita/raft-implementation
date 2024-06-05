package data

import (
	"sync"
	"time"
)

type SafeTimeout struct {
	Value time.Duration
	Mu    sync.Mutex
}

func NewSafeTimeout() *SafeTimeout {
	return &SafeTimeout{Value: 0}
}
