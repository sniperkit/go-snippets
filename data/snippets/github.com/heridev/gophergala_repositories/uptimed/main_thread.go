package main

import "runtime"

var mainThread = make(chan func())

func runOnMainThread(f func()) {
	mainThread <- f
}

func startMainThread() {
	for f := range mainThread {
		go func() {
			runtime.LockOSThread()
			f()
		}()
	}
}
