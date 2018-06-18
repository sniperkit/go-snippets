package binding


/*
#cgo CXXFLAGS: -std=c++0x -Wall -fno-strict-aliasing -I.
#cgo LDFLAGS: -L. -Lcore/capi -L/home/delta54/Qt5.5.0/5.5/gcc/lib -lQt5Core -lQt5Gui -lQt5Widgets -lstdc++  -lcapi

#include "capi.h"

*/
import "C"

import "unsafe"

type SystemTray struct {
    ptr unsafe.Pointer
}

func NewSystemTray() SystemTray {
    var instance SystemTray
    instance.ptr = C.NewSystemTray()
    return instance
}

func (t SystemTray) Init() {
    C.Init(t.ptr)
}

func (t SystemTray) SetIcon(img string) {
    C.SetTrayIcon(t.ptr, C.CString(img))
}

func (t SystemTray) SetToolTip(tooltip string) {
    C.SetTrayToolTip(t.ptr, C.CString(tooltip))
}