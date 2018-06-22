package models

import (
	"time"
)

var (
	defaultInterval = 5 * time.Second
)

type Semaphore struct {
	Inputs      []Link
	ActiveInput *Link
	Interval    time.Duration
	Status      chan SemRequest
	Paused      bool
}

type SemRequest struct {
	Status bool
	Allow  string //link active's street
}

func defaultSemaphore() Semaphore {
	return Semaphore{Inputs: make([]Link, 0), ActiveInput: nil, Interval: defaultInterval, Status: make(chan SemRequest, 1), Paused: false}
}

func (sem *Semaphore) Start() {
	change := time.After(sem.Interval)
	var current int
	if len(sem.Inputs) == 0 {
		return
	}
	sem.ActiveInput = &sem.Inputs[current]

	for {
		select {
		case <-change:
			if !sem.Paused {
				current++
				if current == len(sem.Inputs) {
					current = 0
				}
				sem.ActiveInput = &sem.Inputs[current]
			}
			change = time.After(sem.Interval)
		case req := <-sem.Status:
			sem.Paused = req.Status
			if req.Status {
				for i := 0; i < len(sem.Inputs); i++ {
					if sem.Inputs[i].Name == req.Allow {
						sem.ActiveInput = &sem.Inputs[i]
					}
				}
			} else {
				change = time.After(1 * time.Second)
			}
		}
	}
}
