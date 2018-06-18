package channel

import "time"

func AsyncSend(ch chan<- struct{}) bool {
	select {
	case ch <- struct{}{}:
		return true
	default:
	}
	return false
}

func AsyncSendWithRetry(ch chan<- struct{}, retry int, dur time.Duration) bool {
	for i := 0; i < retry; i++ {
		if AsyncSend(ch) {
			return true
		}
		if dur > 0 {
			time.Sleep(dur)
		}
	}
	return false
}

func AsyncSendError(ch chan<- error, err error) bool {
	select {
	case ch <- err:
		return true
	default:
	}
	return false
}

func AsyncSendErrorWithRetry(ch chan<- error, err error, retry int, dur time.Duration) bool {
	for i := 0; i < retry; i++ {
		if AsyncSendError(ch, err) {
			return true
		}
		if dur > 0 {
			time.Sleep(dur)
		}
	}
	return false
}

func SyncSend(ch chan<- struct{}) {
	ch <- struct{}{}
}

func SyncSendWithTimeout(ch chan<- struct{}, dur time.Duration) bool {
	if dur <= 0 {
		panic("SyncSendWithTimeout dur must greater than zero")
	}

	select {
	case ch <- struct{}{}:
		return true
	case <-time.After(dur):
		return false
	}
}

func SyncSendWithRetry(ch chan<- struct{}, retry int, dur time.Duration) bool {
	if dur <= 0 {
		panic("SyncSendWithRetry dur must greater than zero")
	}

	for i := 0; i < retry; i++ {
		if SyncSendWithTimeout(ch, dur) {
			return true
		}
	}
	return false
}
