package entry

// Panopticon, Copyright © 2015, Huck Ridge Software LLC
// All rights reserved.

// vim:sw=4:ts=4

import "time"

type Entry struct {
	Time    time.Time
	WasIdle bool
	Idle    time.Duration
	// App     string // Don't know how to figure this out yet.
	Title string
}
