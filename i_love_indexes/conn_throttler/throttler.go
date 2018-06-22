package conn_throttler

import "github.com/abiosoft/semaphore"

const (
	MaxConnectionsPerHost = 64
)

var (
	requestsThrottler = map[string]*semaphore.Semaphore{}
)

func Acquire(host string) *semaphore.Semaphore {
	sem, ok := requestsThrottler[host]
	if !ok {
		requestsThrottler[host] = semaphore.New(MaxConnectionsPerHost)
		sem = requestsThrottler[host]
	}
	sem.Acquire()
	return sem
}
