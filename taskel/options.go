package taskel

import (
	"log"
	"syscall"
	"time"

	"github.com/GeorgeMac/tools/wrappers"
)

type option func(s *Scheduler)

// SetTimeout uses standard time.Ticker functionality
func SetTimeout(dur time.Duration) option {
	return func(s *Scheduler) {
		s.timeout = func() Trigger { return wrappers.StructTick(dur) }
	}
}

// SetTimeoutOn allows you to pass a signalling
// channel of your own.
func SetTimeoutOn(t Trigger) option {
	return func(s *Scheduler) {
		s.timeout = func() Trigger {
			return t
		}
	}
}

// SetTimeBetween allows you to define a wait period
// betwen running tasks.
func SetTimeBetween(dur time.Duration) option {
	return func(s *Scheduler) {
		s.between = func() Trigger { return wrappers.StructTick(dur) }
	}
}

// SetTimeBetweenOn allows you to pass a signalling
// channel of your own.
func SetTimeBetweenOn(t Trigger) option {
	return func(s *Scheduler) {
		s.between = func() Trigger {
			return t
		}
	}
}

// SetTerminateOn will trigger termination when it recvs
// on Trigger t.
func SetTerminateOn(t Trigger) option {
	return func(s *Scheduler) {
		s.term = t
	}
}

// SetTerminateOnSignal will trigger termination on the capture
// of a syscall.SIGTERM.
func SetTerminateOnSignal(s *Scheduler) {
	ch := make(chan struct{})
	wrappers.Notify(ch, syscall.SIGTERM)
	s.term = ch
}

// SetLogger enables logging by the taskel scheduler
func SetLogger(logger *log.Logger) option {
	return func(s *Scheduler) {
		s.logger = logger
	}
}
