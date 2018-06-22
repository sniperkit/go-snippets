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
 *     Initial: 2017/08/18        Feng Yifei
 */

package scheduler

import (
	"errors"
	"runtime"
	"time"
)

var (
	// ErrScheduleTimeout happens when task schedule failed during the specific interval.
	ErrScheduleTimeout = errors.New("schedule not available currently")
)

// Pool caches tasks and schedule tasks to work.
type Pool struct {
	queue   chan Task
	workers chan chan Task
}

// New a goroutine pool.
func New(qsize, wsize int) *Pool {
	if wsize == 0 {
		wsize = runtime.NumCPU()
	}

	if qsize < wsize {
		qsize = wsize
	}

	pool := &Pool{
		queue:   make(chan Task, qsize),
		workers: make(chan chan Task, wsize),
	}

	go pool.start()

	for i := 0; i < wsize; i++ {
		StartWorker(pool)
	}

	return pool
}

// Starts the scheduling.
func (p *Pool) start() {
	for {
		select {
		case worker := <-p.workers:
			task := <-p.queue
			worker <- task
		}
	}
}

// Schedule push a task on queue.
func (p *Pool) Schedule(task Task) {
	p.queue <- task
}

// ScheduleWithTimeout try to push a task on queue, if timeout, return false.
func (p *Pool) ScheduleWithTimeout(timeout time.Duration, task Task) error {
	timer := time.NewTimer(timeout)

	select {
	case p.queue <- task:
		timer.Stop()
		return nil
	case <-timer.C:
		return ErrScheduleTimeout
	}
}
