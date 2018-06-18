package binding


/*
#cgo CXXFLAGS: -std=c++0x -Wall -fno-strict-aliasing -I.
#cgo LDFLAGS: -L. -Lcore/capi -L/usr/lib -lQt5Core -lQt5Gui -lQt5Widgets -lstdc++  -lcapi

#include "capi.h"

*/
import "C"

type SystemTray struct {
    ptr C.SystemTray_
}

func NewSystemTray() SystemTray {
    var inst SystemTray;
    inst.ptr = C.NewSystemTray();
    return inst;
}

func (i SystemTray) Init() {
    C.Init(i.ptr)
}

func (i SystemTray) SetTrayIcon(img string) {
    C.SetTrayIcon(i.ptr, C.CString(img))
}

func (i SystemTray) SetTrayToolTip(title string) {
    C.SetTrayToolTip(i.ptr, C.CString(title))
}