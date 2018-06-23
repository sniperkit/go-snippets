package tfcgo

// #include "say.h"
// #include <stdlib.h>
// #include "tensorflow/c/c_api.h"
// #cgo CFLAGS: -I/usr/local/include
// #cgo LDFLAGS: -ltensorflow
import "C"
import "unsafe"

//SaySomething
func SaySomething(saying string) string {
	cs := C.CString(saying)
	say := C.GoString(C.TFCGO_SaySomething(cs))
	C.free(unsafe.Pointer(cs))
	return say
}
