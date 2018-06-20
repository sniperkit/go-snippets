package main

// #include <stdlib.h>
// #include "tensorflow/c/c_api.h"
// #cgo CFLAGS: -I/usr/local/include
// #cgo linux LDFLAGS: -ltensorflow
import "C"

import (
	"unsafe"

	"fmt"

	"github.com/golang/protobuf/proto"
	pb "github.com/tensorflow/tensorflow/tensorflow/go/pb/tensorflow/core/framework"
)

func main() {
	ops, err := registeredOps()
	if err != nil {
		return
	}
	for _, op := range ops.Op {
		fmt.Println(op.Name)
	}
}

func registeredOps() (*pb.OpList, error) {
	buf := C.TF_GetAllOpList()
	defer C.TF_DeleteBuffer(buf)
	var (
		list = new(pb.OpList)
		size = int(buf.length)
		// See: https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices
		data = (*[1 << 30]byte)(unsafe.Pointer(buf.data))[:size:size]
		err  = proto.Unmarshal(data, list)
	)
	return list, err
}
