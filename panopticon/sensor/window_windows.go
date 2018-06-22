package main

// Panopticon, Copyright © 2015, Huck Ridge Software LLC
// All rights reserved.

// vim:ts=4

import (
	"errors"
	"syscall"
	"unsafe"
)

type HWND uintptr
type DWORD uint32
type TickCount DWORD
type LastInputInfo struct {
	cbSize uint32
	dwTime TickCount
}

// Might be wrong.  Not tested.
type MouseMovePoint struct {
	X, Y        int
	Time        DWORD // should be TickCount?
	dwExtraInfo *uint32
}

type LONG int32
type Point struct {
	X, Y LONG
}

var (
	user32                  = syscall.MustLoadDLL("user32.dll")
	getForegroundWindow_W32 = user32.MustFindProc("GetForegroundWindow")
	getWindowText_W32       = user32.MustFindProc("GetWindowTextW")
	getLastInputInfo_W32    = user32.MustFindProc("GetLastInputInfo")
	getCursorPos_W32        = user32.MustFindProc("GetCursorPos")

	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	getTickCount_W32 = kernel32.MustFindProc("GetTickCount")
)
var testHandle HWND

func getForegroundWindow() HWND {
	if testHandle != 0 {
		return testHandle
	}
	windowHandle, _, err := getForegroundWindow_W32.Call()
	if windowHandle == 0 {
		panic("error getting foreground window handle: " + err.Error())
	}
	// log.Printf("windowHandle is %v\n", windowHandle)
	return HWND(windowHandle)
}

func WindowTitle() string {
	// or you can handle the errors in the above if you want to provide some alternative
	windowHandle := getForegroundWindow()
	var buffer [256]uint16
	windowTitleLen, _, _ := getWindowText_W32.Call(uintptr(windowHandle),
		uintptr(unsafe.Pointer(&buffer)), uintptr(256))
	if windowTitleLen == 0 {
		return ""
	}
	return syscall.UTF16ToString(buffer[:windowTitleLen])
}

func GetLastInputInfo() (TickCount, error) {
	lastInputInfo := LastInputInfo{}
	lastInputInfo.cbSize = uint32(unsafe.Sizeof(lastInputInfo))
	rc, _, _ := getLastInputInfo_W32.Call(uintptr(unsafe.Pointer(&lastInputInfo)))
	if int32(rc) == 0 {
		return 0, errors.New("No time returned")
	}
	return lastInputInfo.dwTime, nil
}

func GetCursorPos() (*Point, error) {
	res := Point{}
	b, _, _ := getCursorPos_W32.Call(uintptr(unsafe.Pointer(&res)))
	if int32(b) == 0 {
		return nil, errors.New("No mouse pos available")
	}
	return &res, nil
}

func GetTickCount() TickCount {
	ticks, _, _ := getTickCount_W32.Call()
	return TickCount(ticks)
}
