package tfcgo

// #include <stdlib.h>
// #include "tensorflow/c/c_api.h"
// #cgo CFLAGS: -I/usr/local/include
// #cgo LDFLAGS: -ltensorflow
import "C"

import (
	"strconv"
	"unsafe"

	"github.com/golang/protobuf/proto"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
	pb "github.com/tensorflow/tensorflow/tensorflow/go/pb/tensorflow/core/framework"
)

var constCounter int

//RegisteredOps returns all of the supported TF c/c++ api
func RegisteredOps() (*pb.OpList, error) {
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

func MakeTensorAndOutput(s *op.Scope, value interface{}) (*tf.Tensor, tf.Output, error) {
	t, err := tf.NewTensor(value)
	if err != nil {
		return nil, tf.Output{}, err
	}
	constCounter++
	counter := strconv.Itoa(constCounter)
	op := s.AddOperation(tf.OpSpec{
		Type: "Const",
		Name: "Name_" + counter,
		Attrs: map[string]interface{}{
			"dtype": t.DataType(),
			"value": t,
		},
	})
	return t, op.Output(0), nil
}

func MakeConst(s *op.Scope, value interface{}) (tf.Output, error) {
	t, err := tf.NewTensor(value)
	if err != nil {
		return tf.Output{}, err
	}
	constCounter++
	counter := strconv.Itoa(constCounter)
	op := s.AddOperation(tf.OpSpec{
		Type: "Const",
		Name: "Name_" + counter,
		Attrs: map[string]interface{}{
			"dtype": t.DataType(),
			"value": t,
		},
	})
	return op.Output(0), nil
}
