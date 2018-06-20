package main

import (
	"fmt"
	"runtime"

	"github.com/ctava/tfcgo"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

func main() {

	fmt.Println("go.Version = ", runtime.Version())
	fmt.Println("tf.Version = ", tf.Version())
	fmt.Println("tfcgo.Version = ", tfcgo.Version())

}
