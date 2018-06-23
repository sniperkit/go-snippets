package main

import (
	"fmt"

	"github.com/ctava/tfcgo"
)

func main() {

	ops, err := tfcgo.RegisteredOps()
	if err != nil {
		return
	}
	for _, op := range ops.Op {
		fmt.Println(op.Name)
	}
}
