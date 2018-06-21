package connect

import (
	"runtime"
	"time"
)

var (
	connectSendTimeout   time.Duration
	connectReconnectTime time.Duration
)

func init() {
	connectReconnectTime = time.Millisecond * 100
	connectSendTimeout = time.Millisecond
}

var (
	defaultPoolSize int
)

func init() {
	defaultPoolSize = runtime.NumCPU()
}
