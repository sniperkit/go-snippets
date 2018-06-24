package main

import (
	"os"
	"fmt"
)

func main() {
	fatalNoErr("out should not be used")

}

func log(doing string) {
	fmt.Fprintln(os.Stderr, doing)
}

func fatalNoErr(doing string) {
	log(doing)
	os.Exit(1)
}
