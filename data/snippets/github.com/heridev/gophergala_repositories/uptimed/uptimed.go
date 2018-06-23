package main

/*
#cgo LDFLAGS: -framework Cocoa -framework ApplicationServices -framework IOKit
#cgo CFLAGS: -x objective-c
#import "header.h"
*/
import "C"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"unsafe"
)

var (
	idlePollFrequency = flag.Duration("f", 10*time.Second, "sets idle time polling frequency")
	idleMinDuration   = flag.Duration("m", 5*time.Minute, "minimum idle time to subtract from uptime")
	debug             = flag.Bool("d", false, "print debug info")
	help              = flag.Bool("h", false, "show help message")

	startAt  time.Time     // reference time for when computer started or woke last
	onTime   time.Duration // total time on
	idleTime time.Duration // total idle time
	upTime   time.Duration // total active time
)

func main() {
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// gather boot and wake times
	boot, err := bootTime()
	if err != nil {
		log.Fatal(err)
	}
	wake, err := wakeTime()
	if err != nil {
		log.Fatal(err)
	}

	// used either boot or wake time, whichever is later,
	// as reference starting time
	if wake.After(*boot) {
		startAt = *wake
	} else {
		startAt = *boot
	}

	// start main thread loop
	go startMainThread()
	// start OSX app
	runOnMainThread(func() { C.StartApp() })

	// start idle poller
	go startSysIdleTimePoller(*idlePollFrequency, *idleMinDuration)
	// update ticker
	updateTicker := time.NewTicker(time.Second).C

	// update the status bar with current upTime
	update := func() {
		onTime = time.Since(startAt)
		upTime = onTime - idleTime
		setMenuLabel(formatDuration(upTime))
		if *debug {
			fmt.Printf("onTime: %q\tidleTime: %q\tupTime: %q\n", onTime, idleTime, upTime)
		}
	}
	for {
		select {
		case <-updateTicker:
			update()
		case idleInc := <-idleTicker:
			// idle time has increased
			idleTime += idleInc
			update()
		}
	}
}

// setMenuLabel changes the text in the app's menu bar
func setMenuLabel(l string) {
	runOnMainThread(func() {
		cs := C.CString(l)
		C.SetLabelText(cs)
		C.free(unsafe.Pointer(cs))
	})
}
