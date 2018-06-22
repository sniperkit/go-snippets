package basic

import (
	"log"
	"testing"
)

type Code int

const (
	CodeNormal Code = iota + 900
	CodeErr

	CodeTimeout
)

func TestIota(t *testing.T) {

	log.Printf("%v ==> %T", CodeErr, CodeErr)

	log.Printf("%v ==> %T", CodeTimeout, CodeTimeout)
}
