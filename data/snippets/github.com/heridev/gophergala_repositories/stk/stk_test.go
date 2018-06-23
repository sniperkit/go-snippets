package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestStkErr(t *testing.T) {
	//Log automatically writes to Stderr
	//If not using official log, simply pass in os.Stderr as an io.Writer
	//to make it work
	fmt.Println("Const os.Stderr", os.Stderr)
	fmt.Println("Hey", log.Flags())
}
