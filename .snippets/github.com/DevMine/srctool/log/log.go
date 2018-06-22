// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package log provides simple logging routines.
//
// Everything is logged into stderr in order to keep stdout empty. This is
// required because the final JSON must be output into stdout by default.
package log

import (
	"fmt"
	"os"
	"syscall"

	"github.com/gilliek/go-xterm256/xterm256"
	"golang.org/x/crypto/ssh/terminal"
)

var debugMode = false

// SetDebugMode allows to enable/disable debug mode.
func SetDebugMode(val bool) {
	debugMode = val
}

// Success prints success messages.
func Success(a ...interface{}) {
	write(xterm256.Green, "success", fmt.Sprint(a...))
}

// Info prints info messages.
func Info(a ...interface{}) {
	write(xterm256.Blue, "info", fmt.Sprint(a...))
}

// Debug prints debug messages only if the debug mode is enabled.
func Debug(a ...interface{}) {
	if debugMode {
		write(xterm256.White, "debug", fmt.Sprint(a...))
	}
}

// Fail prints error messages without exiting the program.
func Fail(a ...interface{}) {
	write(xterm256.Red, "error", fmt.Sprint(a...))
}

// Fatal prints error messages, then exits the program with status code 1.
func Fatal(a ...interface{}) {
	write(xterm256.Red, "fatal", fmt.Sprint(a...))
	os.Exit(1)
}

// write is the low level routine for printing log messages with color c and label l.
func write(c xterm256.Color, l, msg string) {
	label := l
	if terminal.IsTerminal(syscall.Stderr) {
		label = xterm256.Sprint(c, l)
	}

	fmt.Fprintf(os.Stderr, "[%-7s] %s\n", label, msg)
}
