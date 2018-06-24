package tsm

/*
#include <stdlib.h>
#include <libtsm.h>

extern void vteWrite(struct tsm_vte *vte, char *u8, size_t len, void *data);
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type VTEWriteFunc func(b []byte)

//export vteWrite
func vteWrite(vte *C.struct_tsm_vte, u8 *C.char, len C.size_t, data unsafe.Pointer) {
	f := *(*VTEWriteFunc)(data)
	b := C.GoBytes(unsafe.Pointer(u8), C.int(len))
	f(b)
}

type VTEModifiers uint

const (
	VTEModifierShift VTEModifiers = 1 << iota
	VTEModifierLock
	VTEModifierControl
	VTEModifierAlt
	VTEModifierLogo
)

type VTE struct {
	vte *C.struct_tsm_vte

	Screen *Screen
}

func NewVTE(s *Screen, f VTEWriteFunc) *VTE {
	vte := &VTE{Screen: s}
	ret := C.tsm_vte_new(&vte.vte, s.s, (*[0]byte)(C.vteWrite), unsafe.Pointer(&f), nil, nil)
	if ret != 0 {
		panic("tsm: failed to create VTE")
	}
	runtime.SetFinalizer(s, func(vte *VTE) {
		C.tsm_vte_unref(vte.vte)
	})
	return vte
}

func (vte *VTE) SetPalette(palette string) {
	cpalette := C.CString(palette)
	C.tsm_vte_set_palette(vte.vte, cpalette)
	C.free(unsafe.Pointer(cpalette))
}

func (vte *VTE) Reset() {
	C.tsm_vte_reset(vte.vte)
}

func (vte *VTE) HardReset() {
	C.tsm_vte_hard_reset(vte.vte)
}

func (vte *VTE) Input(b []byte) {
	cb := C.CBytes(b)
	C.tsm_vte_input(vte.vte, (*C.char)(cb), C.size_t(len(b)))
	C.free(cb)
}

func (vte *VTE) HandleKeyboard(keysym, ascii uint32, mods VTEModifiers, unicode rune) bool {
	return bool(C.tsm_vte_handle_keyboard(vte.vte, C.uint32_t(keysym), C.uint32_t(ascii), C.uint(mods), C.uint32_t(unicode)))
}
