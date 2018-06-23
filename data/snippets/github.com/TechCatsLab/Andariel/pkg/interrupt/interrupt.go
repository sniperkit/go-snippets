/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/03/29        From Kubernetes(Google)
 */

package interrupt

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var terminationSignals = []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGALRM}

// Handler guarantees execution of notifications after a critical section (the function passed
// to a Run method), even in the presence of process termination. It guarantees exactly once
// invocation of the provided notify functions.
type Handler struct {
	notify []func()
	final  func(os.Signal)
	once   sync.Once
}

// Chain creates a new handler that invokes all notify functions when the critical section exits
// and then invokes the optional handler's notifications. This allows critical sections to be
// nested without losing exactly once invocations. Notify functions can invoke any cleanup needed
// but should not exit (which is the responsibility of the parent handler).
func Chain(handler *Handler, notify ...func()) *Handler {
	if handler == nil {
		return New(nil, notify...)
	}
	return New(handler.Signal, append(notify, handler.Close)...)
}

// New creates a new handler that guarantees all notify functions are run after the critical
// section exits (or is interrupted by the OS), then invokes the final handler. If no final
// handler is specified, the default final is `os.Exit(1)`. A handler can only be used for
// one critical section.
func New(final func(os.Signal), notify ...func()) *Handler {
	return &Handler{
		final:  final,
		notify: notify,
	}
}

// Close executes all the notification handlers if they have not yet been executed.
func (h *Handler) Close() {
	h.once.Do(func() {
		for _, fn := range h.notify {
			fn()
		}
	})
}

// Signal is called when an os.Signal is received, and guarantees that all notifications
// are executed, then the final handler is executed. This function should only be called once
// per Handler instance.
func (h *Handler) Signal(s os.Signal) {
	h.once.Do(func() {
		for _, fn := range h.notify {
			fn()
		}
		if h.final == nil {
			os.Exit(1)
		}
		h.final(s)
	})
}

// Run ensures that any notifications are invoked after the provided fn exits (even if the
// process is interrupted by an OS termination signal). Notifications are only invoked once
// per Handler instance, so calling Run more than once will not behave as the user expects.
func (h *Handler) Run(fn func() error) error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, terminationSignals...)
	defer func() {
		signal.Stop(ch)
		close(ch)
	}()
	go func() {
		sig, ok := <-ch
		if !ok {
			return
		}
		h.Signal(sig)
	}()
	defer h.Close()
	return fn()
}

// 等待直到信号量出现
func (h *Handler) Wait() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, terminationSignals...)

	for {
		sig, ok := <-ch
		if !ok {
			return
		}

		if h.final != nil {
			h.final(sig)
		}

		if sig != syscall.SIGALRM {
			signal.Stop(ch)
			close(ch)
		}
	}
}
