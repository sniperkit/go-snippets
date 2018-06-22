package main

// Panopticon, Copyright © 2015, Huck Ridge Software LLC
// All rights reserved.

// vim:sw=4:ts=4

import (
	"log"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWindowTitle(t *testing.T) {
	// This is the empirically determined handle of a mintty in $HOME.  This
	// would have to be updated every time said mintty is restarted.
	testHandle = 3278308
	expected := "~"

	Convey("Get the current window title", t, func() {
		title := WindowTitle()
		Convey("Titles should match", func() {
			So(title, ShouldEqual, expected)
		})
		Convey("Lengths should match", func() {
			So(len(title), ShouldEqual, len(expected))
		})
	})
}

// Not sure how to test this.  I guess, just make sure it's non-zero?
func TestGetLastInputInfo(t *testing.T) {
	Convey("Get tick count of the most recent keyboard activity", t, func() {
		if ticks, err := GetLastInputInfo(); err != nil {
			So(err, ShouldBeNil) // obviously an automatic failure
		} else {
			Convey("tick count is non-zero", func() {
				log.Printf("GetLastInputInfo is %v", ticks)
				So(ticks, ShouldBeGreaterThan, 0)
			})
		}
	})
}

// Not sure how to test this, either.  I guess, just make sure it's non-zero,
// too?  Which of course will fail if you run the test with the mouse in the
// upper left corner of the screen.
func TestGetCursorPos(t *testing.T) {
	Convey("Get the mouse position", t, func() {
		if pos, err := GetCursorPos(); err != nil {
			So(err, ShouldBeNil) // obviously an automatic failure
		} else {
			// log.Printf("Mouse pos is %v", pos)
			Convey("Mouse pos x is non-zero", func() {
				So(pos.X, ShouldBeGreaterThan, 0)
			})
			Convey("Mouse pos y is non-zero", func() {
				So(pos.Y, ShouldBeGreaterThan, 0)
			})
		}
	})
}

func TestMakeEntry(t *testing.T) {
	testHandle = 3278308
	Convey("Make an entry", t, func() {
		e, err := MakeEntry()
		log.Printf("e is %v", *e)
		Convey("Got an entry", func() {
			So(err, ShouldBeNil)
		})
		if err == nil {
			// There's a race condition here: the returned Entry uses time.Now(); the
			// time.Now() *now* is going to be later.  Mitigate (or totally avoid?) this
			// by rounding to the nearest second.
			Convey("times match", func() {
				t1 := e.Time.Round(time.Second)
				t2 := time.Now().Round(time.Second)
				// Comparing them directly didn't work even when they were equal, so stringify them.
				So(t1.String(), ShouldEqual, t2.String())
			})
			Convey("Titles match", func() {
				So(e.Title, ShouldEqual, "~")
			})
		}
	})
}

func TestGetTickCount(t *testing.T) {
	Convey("Get the tick count", t, func() {
		tc := GetTickCount()
		So(tc, ShouldBeGreaterThan, 0)
	})
}
