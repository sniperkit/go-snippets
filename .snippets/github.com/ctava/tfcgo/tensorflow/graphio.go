package tensorflow

// #include <stdlib.h>
// #include "tensorflow/c/c_api.h"
// #cgo CFLAGS: -I/usr/local/include
// #cgo LDFLAGS: -ltensorflow
import "C"

import (
	"bufio"
	"log"
	"os"
	"unsafe"

	"github.com/golang/protobuf/proto"
	pb "github.com/tensorflow/tensorflow/tensorflow/go/pb/tensorflow/core/framework"
)

//WriteGraphAsText writes out a graph in text
func (g *Graph) WriteGraphAsText() error {

	graphDefBuf := C.TF_NewBuffer()
	defer C.TF_DeleteBuffer(graphDefBuf)
	status := C.TF_NewStatus()
	C.TF_GraphToGraphDef(g.c, graphDefBuf, status)
	// if err := status.Err(); err != nil { //TODO
	// 	return err
	// }

	var (
		graphDef        = new(pb.GraphDef)
		graphDefBufSize = int(graphDefBuf.length)
		graphDefData    = (*[1 << 30]byte)(unsafe.Pointer(graphDefBuf.data))[:graphDefBufSize:graphDefBufSize]
		err             = proto.Unmarshal(graphDefData, graphDef)
	)

	file, err := os.Create("graph.pbtxt")
	defer file.Close()
	w := bufio.NewWriter(file)
	defer w.Flush()

	graphText := graphDef.String()
	_, err = w.Write([]byte(graphText))
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
