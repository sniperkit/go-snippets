package tsm

/*
#include <libtsm.h>

extern void get_screen_attr_bitfield(struct tsm_screen_attr *attr, bool *bold,
	bool *underline, bool *inverse, bool *protect, bool *blink);
extern void set_screen_attr_bitfield(struct tsm_screen_attr *attr, bool bold,
	bool underline, bool inverse, bool protect, bool blink);

extern int screenDraw(struct tsm_screen *con, uint32_t id, uint32_t *ch,
	size_t len, unsigned int width, unsigned int posx, unsigned int posy,
	struct tsm_screen_attr *attr, tsm_age_t age, void *data);
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type ScreenAttr struct {
	FCCode, BCCode int8 // foreground, background color codes (<0 for rgb)
	FR, FG, FB, BR, BG, BB uint8 // foreground/background red/green/blue
	Bold, Underline, Inverse, Protect, Blink bool
}

func screenAttrFromC(cattr *C.struct_tsm_screen_attr) *ScreenAttr {
	var bold, underline, inverse, protect, blink C.bool
	C.get_screen_attr_bitfield(cattr, &bold, &underline, &inverse, &protect, &blink)

	return &ScreenAttr{
		FCCode: int8(cattr.fccode),
		BCCode: int8(cattr.bccode),
		FR: uint8(cattr.fr),
		FG: uint8(cattr.fg),
		FB: uint8(cattr.fb),
		BR: uint8(cattr.br),
		BG: uint8(cattr.bg),
		BB: uint8(cattr.bb),
		Bold: bool(bold),
		Underline: bool(underline),
		Inverse: bool(inverse),
		Protect: bool(protect),
		Blink: bool(blink),
	}
}

func (attr *ScreenAttr) toC() *C.struct_tsm_screen_attr {
	cattr := C.struct_tsm_screen_attr{
		fccode: C.int8_t(attr.FCCode),
		bccode: C.int8_t(attr.BCCode),
		fr: C.uint8_t(attr.FR),
		fg: C.uint8_t(attr.FG),
		fb: C.uint8_t(attr.FB),
		br: C.uint8_t(attr.BR),
		bg: C.uint8_t(attr.BG),
		bb: C.uint8_t(attr.BB),
	}
	C.set_screen_attr_bitfield(&cattr, C.bool(attr.Bold), C.bool(attr.Underline), C.bool(attr.Inverse), C.bool(attr.Protect), C.bool(attr.Blink))
	return &cattr
}

type ScreenDrawFunc func(id uint32, s string, width, posx, posy uint, attr *ScreenAttr, age uint32) bool

//export screenDraw
func screenDraw(s *C.struct_tsm_screen, id C.uint32_t, ch *C.uint32_t, len C.size_t, width, posx, posy C.uint, cattr *C.struct_tsm_screen_attr, age C.tsm_age_t, data unsafe.Pointer) C.int {
	f := *(*ScreenDrawFunc)(data)
	str := C.GoStringN((*C.char)(unsafe.Pointer(ch)), C.int(len))
	attr := screenAttrFromC(cattr)
	ok := f(uint32(id), str, uint(width), uint(posx), uint(posy), attr, uint32(age))
	if !ok {
		return 1
	}
	return 0
}

type ScreenFlags uint

const (
	ScreenInsertMode ScreenFlags = 1 << iota
	ScreenAutoWrap
	ScreenRelOrigin
	ScreenInverse
	ScreenHideCursor
	ScreenFixedPos
	ScreenAlternate
)

type Screen struct {
	s *C.struct_tsm_screen
}

func NewScreen() *Screen {
	s := &Screen{}
	if ret := C.tsm_screen_new(&s.s, nil, nil); ret != 0 {
		panic("tsm: failed to create screen")
	}
	runtime.SetFinalizer(s, func(s *Screen) {
		C.tsm_screen_unref(s.s)
	})
	return s
}

func (s *Screen) Width() uint {
	return uint(C.tsm_screen_get_width(s.s))
}

func (s *Screen) Height() uint {
	return uint(C.tsm_screen_get_height(s.s))
}

func (s *Screen) Resize(w, h uint) {
	if ret := C.tsm_screen_resize(s.s, C.uint(w), C.uint(h)); ret != 0 {
		panic("tsm: failed to resize screen")
	}
}

func (s *Screen) SetMargins(top, bottom uint) {
	if ret := C.tsm_screen_resize(s.s, C.uint(top), C.uint(bottom)); ret != 0 {
		panic("tsm: failed to set screen margins")
	}
}

func (s *Screen) SetMaxScrollback(max uint) {
	C.tsm_screen_set_max_sb(s.s, C.uint(max))
}

func (s *Screen) ClearScrollback() {
	C.tsm_screen_clear_sb(s.s)
}

func (s *Screen) ScrollbackUp(num uint) {
	C.tsm_screen_sb_up(s.s, C.uint(num))
}

func (s *Screen) ScrollbackDown(num uint) {
	C.tsm_screen_sb_down(s.s, C.uint(num))
}

func (s *Screen) ScrollbackPageUp(num uint) {
	C.tsm_screen_sb_page_up(s.s, C.uint(num))
}

func (s *Screen) ScrollbackPageDown(num uint) {
	C.tsm_screen_sb_page_down(s.s, C.uint(num))
}

func (s *Screen) ScrollbackReset() {
	C.tsm_screen_sb_reset(s.s)
}

func (s *Screen) SetDefAttr(attr *ScreenAttr) {
	C.tsm_screen_set_def_attr(s.s, attr.toC())
}

func (s *Screen) Reset() {
	C.tsm_screen_reset(s.s)
}

func (s *Screen) SetFlags(flags ScreenFlags) {
	C.tsm_screen_set_flags(s.s, C.uint(flags))
}

func (s *Screen) ResetFlags(flags ScreenFlags) {
	C.tsm_screen_reset_flags(s.s, C.uint(flags))
}

func (s *Screen) Flags() ScreenFlags {
	return ScreenFlags(C.tsm_screen_get_flags(s.s))
}

func (s *Screen) CursorX() uint {
	return uint(C.tsm_screen_get_cursor_x(s.s))
}

func (s *Screen) CursorY() uint {
	return uint(C.tsm_screen_get_cursor_y(s.s))
}

func (s *Screen) SetTabstop() {
	C.tsm_screen_set_tabstop(s.s)
}

func (s *Screen) ResetTabstop() {
	C.tsm_screen_reset_tabstop(s.s)
}

func (s *Screen) ResetAllTabstops() {
	C.tsm_screen_reset_all_tabstops(s.s)
}

func (s *Screen) Write(r rune, attr *ScreenAttr) {
	C.tsm_screen_write(s.s, C.tsm_symbol_t(r), attr.toC())
}

func (s *Screen) Newline() {
	C.tsm_screen_newline(s.s)
}

func (s *Screen) ScrollUp(num uint) {
	C.tsm_screen_scroll_up(s.s, C.uint(num))
}

func (s *Screen) ScrollDown(num uint) {
	C.tsm_screen_scroll_down(s.s, C.uint(num))
}

func (s *Screen) MoveTo(x, y uint) {
	C.tsm_screen_move_to(s.s, C.uint(x), C.uint(y))
}

func (s *Screen) MoveUp(num uint, scroll bool) {
	C.tsm_screen_move_up(s.s, C.uint(num), C.bool(scroll))
}

func (s *Screen) MoveDown(num uint, scroll bool) {
	C.tsm_screen_move_down(s.s, C.uint(num), C.bool(scroll))
}

func (s *Screen) MoveLeft(num uint) {
	C.tsm_screen_move_left(s.s, C.uint(num))
}

func (s *Screen) MoveRight(num uint) {
	C.tsm_screen_move_right(s.s, C.uint(num))
}

func (s *Screen) MoveLineEnd() {
	C.tsm_screen_move_line_end(s.s)
}

func (s *Screen) MoveLineHome() {
	C.tsm_screen_move_line_home(s.s)
}

func (s *Screen) TabRight(num uint) {
	C.tsm_screen_tab_right(s.s, C.uint(num))
}

func (s *Screen) TabLeft(num uint) {
	C.tsm_screen_tab_left(s.s, C.uint(num))
}

func (s *Screen) InsertLines(num uint) {
	C.tsm_screen_insert_lines(s.s, C.uint(num))
}

func (s *Screen) DeleteLines(num uint) {
	C.tsm_screen_delete_lines(s.s, C.uint(num))
}

func (s *Screen) InsertChars(num uint) {
	C.tsm_screen_insert_chars(s.s, C.uint(num))
}

func (s *Screen) DeleteChars(num uint) {
	C.tsm_screen_delete_chars(s.s, C.uint(num))
}

func (s *Screen) EraseCursor() {
	C.tsm_screen_erase_cursor(s.s)
}

func (s *Screen) EraseChars(num uint) {
	C.tsm_screen_erase_chars(s.s, C.uint(num))
}

func (s *Screen) EraseCursorToEnd(protect bool) {
	C.tsm_screen_erase_cursor_to_end(s.s, C.bool(protect))
}

func (s *Screen) EraseHomeToCursor(protect bool) {
	C.tsm_screen_erase_home_to_cursor(s.s, C.bool(protect))
}

func (s *Screen) EraseCurrentLine(protect bool) {
	C.tsm_screen_erase_current_line(s.s, C.bool(protect))
}

func (s *Screen) EraseScreenToCursor(protect bool) {
	C.tsm_screen_erase_screen_to_cursor(s.s, C.bool(protect))
}

func (s *Screen) EraseCursorToScreen(protect bool) {
	C.tsm_screen_erase_cursor_to_screen(s.s, C.bool(protect))
}

func (s *Screen) EraseScreen(protect bool) {
	C.tsm_screen_erase_screen(s.s, C.bool(protect))
}

func (s *Screen) SelectionReset() {
	C.tsm_screen_selection_reset(s.s)
}

func (s *Screen) SelectionStart(posx, posy uint) {
	C.tsm_screen_selection_start(s.s, C.uint(posx), C.uint(posy))
}

func (s *Screen) SelectionTarget(posx, posy uint) {
	C.tsm_screen_selection_target(s.s, C.uint(posx), C.uint(posy))
}

func (s *Screen) SelectionCopy() string {
	var sel *C.char
	n := C.tsm_screen_selection_copy(s.s, &sel)
	if int(n) < 0 {
		panic("tsm: failed to copy screen selection")
	}
	return C.GoStringN(sel, n)
}

func (s *Screen) Draw(f ScreenDrawFunc) uint32 {
	age := C.tsm_screen_draw(s.s, (*[0]byte)(C.screenDraw), unsafe.Pointer(&f))
	return uint32(age)
}
