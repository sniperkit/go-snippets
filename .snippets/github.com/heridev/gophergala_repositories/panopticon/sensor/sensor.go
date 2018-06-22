package main

// Panopticon, Copyright © 2015, Huck Ridge Software LLC
// All rights reserved.

// vim:sw=4:ts=4

import (
	"time"

	"github.com/gophergala/panopticon/entry"
)

var prevMousePos Point
var lastMouseMovement = time.Now()
var wasIdle = false

const CONSIDERED_IDLE = 60

func MakeEntry() (*entry.Entry, error) {
	mousePos, err := GetCursorPos()
	if err != nil {
		return nil, err
	}
	kbdLastActive, err := GetLastInputInfo()
	if err != nil {
		return nil, err
	}
	var mouseIdleTime time.Duration
	if prevMousePos == *mousePos {
		mouseIdleTime = time.Now().Sub(lastMouseMovement) / time.Millisecond
	} else {
		lastMouseMovement = time.Now()
		prevMousePos = *mousePos
	}
	kbdIdleTime := time.Duration(int64(GetTickCount()-kbdLastActive) * int64(time.Millisecond))
	idle := mouseIdleTime
	if kbdIdleTime < mouseIdleTime {
		idle = kbdIdleTime
	}
	idleTime := time.Duration(idle * time.Millisecond)
	isIdle := idleTime > CONSIDERED_IDLE
	e := entry.Entry{
		Time:    time.Now(),
		WasIdle: wasIdle,
		Idle:    idleTime,
		Title:   WindowTitle()}
	wasIdle = isIdle
	return &e, nil
}
