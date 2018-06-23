package taskel

import (
	"testing"
	"time"
)

var errmsg string = "Expected %v, Got %v\n"

func Test_Scheduler_New(t *testing.T) {
	var called bool
	fn := func() {
		called = true
	}

	s := New(AfterTaskFunc(fn))
	// check nil was not returned
	if s == nil {
		t.Fatalf(errmsg, "not nil", s)
	}

	// run scheduler task and check it handles notification
	// and function delegation correctly.
	ch := make(chan struct{})
	s.t.Run(ch)

	// ensure called is now set to true
	if !called {
		t.Fatalf(errmsg, true, called)
	}

	// ensure we don't block on notify channel `ch`
	select {
	case <-ch:
		return
	default:
		t.Fatalf(errmsg, "to not block", "blocked")
	}
}

func Test_Scheduler_New_WithTimeoutChannel(t *testing.T) {
	ch := make(chan struct{}, 1)
	s := New(AfterTaskFunc(func() {}), SetTimeoutOn(ch))

	ch <- struct{}{}

	select {
	case <-s.timeout():
		return
	default:
		t.Fatalf(errmsg, "timeout channel to not be block", "timeout channel blocked")
	}
}

func Test_Scheduler_Terminate(t *testing.T) {
	// buffered termination channel
	term := make(chan struct{}, 1)

	var calls int
	s := New(AfterTaskFunc(func() {
		// wait for a second (long task)
		<-time.Tick(time.Second)
		// set called to true on completion of task
		calls++
	}),
		SetTerminateOn(term))

	// begin task
	end := make(chan struct{})
	go s.Begin(end)

	// send on terminate
	close(term)

	// wait for task to be completed
	<-end
	// called should have been incremented
	if calls != 1 {
		t.Errorf(errmsg, 1, calls)
	}
}

func Test_Scheduler_TerminateWithTimeout(t *testing.T) {
	// buffered termination channel
	term, timeout := make(chan struct{}, 1), make(chan struct{}, 1)

	var calls int
	s := New(AfterTaskFunc(func() {
		// block indefinitely
		<-Trigger(nil)
		// set called to true on completion of task
		calls++
	}),
		SetTerminateOn(term),
		SetTimeoutOn(timeout))

	// begin task
	end := make(chan struct{})
	go s.Begin(end)

	// send on term and timeout
	close(term)
	close(timeout)

	// block until Begin closes end
	<-end

	// called should have not finished
	if calls != 0 {
		t.Errorf(errmsg, 0, calls)
	}
}

func Test_Scheduler_TerminateWithBetween(t *testing.T) {
	// buffered termination channel
	term, between := make(chan struct{}, 1), make(chan struct{}, 0)

	var calls int
	s := New(AfterTaskFunc(func() {
		// set called to true on completion of task
		calls++
	}),
		SetTerminateOn(term),
		SetTimeBetweenOn(between))

	// begin task
	end := make(chan struct{})
	go s.Begin(end)

	// Allow for begin to unblock five times
	for i := 0; i < 4; i++ {
		between <- struct{}{}
	}

	// send on term and timeout
	close(term)

	// block until Begin closes end
	<-end

	// should have been called 5 times
	if calls != 5 {
		t.Errorf(errmsg, 5, calls)
	}
}
