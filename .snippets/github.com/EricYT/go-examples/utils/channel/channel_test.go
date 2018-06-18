package channel

import (
	"errors"
	"testing"
	"time"
)

func TestAsyncSend(t *testing.T) {
	ch := make(chan<- struct{})
	if AsyncSend(ch) {
		t.Fatalf("AsyncSend should faild when receiver is not ready")
	}
	ch = make(chan<- struct{}, 1)
	if !AsyncSend(ch) {
		t.Fatalf("AsyncSend should not faild when the channel is buffered")
	}
	if AsyncSend(ch) {
		t.Fatalf("AsyncSend should faild when the channel has not enough buffers")
	}
}

func TestAsyncSendWithRetry(t *testing.T) {
	ch := make(chan<- struct{})
	if AsyncSendWithRetry(ch, 3, time.Second*1) {
		t.Fatalf("AsyncSendWithRetry should faild even though retry three times")
	}
	ch = make(chan<- struct{}, 1)
	if !AsyncSendWithRetry(ch, 3, time.Second*1) {
		t.Fatalf("AsyncSendWithRetry should not faild because there is a buffer")
	}

	doubleCh := make(chan struct{})
	go func() {
		time.Sleep(time.Second * 1)
		<-doubleCh
	}()
	if !AsyncSendWithRetry(doubleCh, 3, time.Second*1) {
		t.Fatalf("AsyncSendWithRetry should succeed after try onece")
	}
}

var (
	ErrorTest error = errors.New("channel: test error")
)

func TestAsyncSendError(t *testing.T) {
	ch := make(chan<- error)
	if AsyncSendError(ch, ErrorTest) {
		t.Fatalf("AsyncSendError should faild when receiver is not ready")
	}
	ch = make(chan<- error, 1)
	if !AsyncSendError(ch, ErrorTest) {
		t.Fatalf("AsyncSendError should not faild when the channel is buffered")
	}
	if AsyncSendError(ch, ErrorTest) {
		t.Fatalf("AsyncSendError should faild when the channel has not enough buffers")
	}
}

func TestAsyncSendErrorWithRetry(t *testing.T) {
	ch := make(chan<- error)
	if AsyncSendErrorWithRetry(ch, ErrorTest, 3, time.Second*1) {
		t.Fatalf("AsyncSendErrorWithRetry should faild even though retry three times")
	}
	ch = make(chan<- error, 1)
	if !AsyncSendErrorWithRetry(ch, ErrorTest, 3, time.Second*1) {
		t.Fatalf("AsyncSendErrorWithRetry should not faild because there is a buffer")
	}

	doubleCh := make(chan error)
	go func() {
		time.Sleep(time.Second * 1)
		<-doubleCh
	}()
	if !AsyncSendErrorWithRetry(doubleCh, ErrorTest, 3, time.Second*1) {
		t.Fatalf("AsyncSendErrorWithRetry should succeed after try onece")
	}
}

func TestSyncSendWithTimeout(t *testing.T) {
	ch := make(chan<- struct{})
	if SyncSendWithTimeout(ch, time.Second*1) {
		t.Fatalf("SyncSendWithTimeout should not succeed when receiver is not ready")
	}
	ch = make(chan<- struct{}, 1)
	if !SyncSendWithTimeout(ch, time.Second*1) {
		t.Fatalf("SyncSendWithTimeout should not failed when channel is buffeded")
	}
}

func TestSyncSendWithRetry(t *testing.T) {
	ch := make(chan struct{})
	go func() {
		time.Sleep(time.Second * 1)
		<-ch
	}()
	if !SyncSendWithRetry(ch, 3, time.Second*1) {
		t.Fatalf("SyncSendWithRetry should succeed after sleep a while")
	}
}
