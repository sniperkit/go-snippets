package demangle

/*
#include <demangle.h>
#include <stdlib.h>
#cgo CFLAGS: -I/usr/include/libiberty
#cgo LDFLAGS: -liberty

static const struct demangler_engine *demangler(int i) {
	return libiberty_demanglers + i;
}
*/
import "C"
import "unsafe"

type Style C.enum_demangling_styles

type Engine struct {
	Name  string
	Style Style
	Doc   string
}

const (
	None    Style = C.no_demangling
	Unknown Style = C.unknown_demangling
	Auto    Style = C.auto_demangling
	Gnu     Style = C.gnu_demangling
	Lucid   Style = C.lucid_demangling
	Arm     Style = C.arm_demangling
	HP      Style = C.hp_demangling
	Edg     Style = C.edg_demangling
	GnuV3   Style = C.gnu_v3_demangling
	Java    Style = C.java_demangling
	Gnat    Style = C.gnat_demangling
	Dlang   Style = C.dlang_demangling
)

const (
	ANSI    = C.DMGL_ANSI
	PARAMS  = C.DMGL_PARAMS
	VERBOSE = C.DMGL_VERBOSE
	TYPES   = C.DMGL_TYPES
)

func Cplus(mangled string, options int) string {
	cmangled := C.CString(mangled)
	defer C.free(unsafe.Pointer(cmangled))
	return C.GoString(C.cplus_demangle(cmangled, C.int(options)))
}

func CplusMangleOpname(opname string, options int) string {
	copname := C.CString(opname)
	defer C.free(unsafe.Pointer(copname))
	return C.GoString(C.cplus_mangle_opname(copname, C.int(options)))
}

func CplusSetStyle(style Style) Style {
	return Style(C.cplus_demangle_set_style(C.enum_demangling_styles(style)))
}

func CplusNameToStyle(name string) Style {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Style(C.cplus_demangle_name_to_style(cname))
}

func CurrentStyle() Style {
	return Style(C.current_demangling_style)
}

var (
	engines []Engine
)

func init() {
	for i := 0; ; i++ {
		p := C.demangler(C.int(i))
		if p.demangling_style == C.unknown_demangling {
			break
		}
		engines = append(engines, Engine{
			Name:  C.GoString(p.demangling_style_name),
			Style: Style(p.demangling_style),
			Doc:   C.GoString(p.demangling_style_doc),
		})
	}
}

func Engines() []Engine {
	return engines
}
