package main

/*
#include "ApplicationServices/ApplicationServices.h"
// retrieve idle time
static inline int64_t systemIdleTime(void) {
  CFTimeInterval timeSinceLastEvent;
  timeSinceLastEvent = CGEventSourceSecondsSinceLastEventType(kCGEventSourceStateCombinedSessionState, kCGAnyInputEventType);
  return timeSinceLastEvent * 1000000000;
}
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"syscall"
	"time"
)

var (
	idleTicker = make(chan time.Duration)
	sleepAt    time.Time
)

// called from objective-c on sleep
//export SleepEvent
func SleepEvent() {
	sleepAt = time.Now()
	idleTicker <- systemIdleTime()
	//fmt.Printf("Sleep event: %q\n", sleepAt)
}

// called from objective-c on wake
//export WakeEvent
func WakeEvent() {
	d := time.Since(sleepAt)
	idleTicker <- d
	//fmt.Printf("Wake event. Slept for %q\n", d)
}

// startSysIdleTimePoller watches the idle time as reported by: systemCGEventSourceSecondsSinceLastEventType(kCGEventSourceStateCombinedSessionState, kCGAnyInputEventType);
// when idle time decreases, report previous idle time if it's greater than the minimum provided
// to be considered idle
func startSysIdleTimePoller(freq time.Duration, min time.Duration) {
	var prev, curr time.Duration
	for _ = range time.NewTicker(freq).C {
		curr = systemIdleTime()
		if curr < prev && prev >= min {
			idleTicker <- (prev + freq)
		}
		prev = curr
	}
}

func systemIdleTime() time.Duration {
	return time.Duration(C.systemIdleTime())
}

func bootTime() (*time.Time, error) {
	return sysCtlTimeByName("kern.boottime")
}

func sleepTime() (*time.Time, error) {
	return sysCtlTimeByName("kern.sleeptime")
}

func wakeTime() (*time.Time, error) {
	return sysCtlTimeByName("kern.waketime")
}

// sysCtlTimeByName reads a sysctl time
func sysCtlTimeByName(name string) (*time.Time, error) {
	v, err := syscall.Sysctl(name)
	if err != nil {
		return nil, fmt.Errorf("%s error: %q", name, err)
	}

	var secs int64
	buf := bytes.NewBufferString(v)
	if err = binary.Read(buf, binary.LittleEndian, &secs); err != nil {
		return nil, fmt.Errorf("binary.Read error: %q", err)
	}
	u := time.Unix(int64(secs), 0)
	return &u, nil
}
