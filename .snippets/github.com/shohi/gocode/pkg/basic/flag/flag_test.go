package flag

import (
	"flag"
	"log"
	"testing"
)

var myStr string

func setup() {
	flag.StringVar(&myStr, "-mystr", "hello", "set str")
}

func TestFlag(t *testing.T) {
	setup()
	flag.Parse()

	log.Println(myStr)
}
