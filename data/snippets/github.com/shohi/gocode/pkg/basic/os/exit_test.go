package os

import (
	"os"
	"testing"
)

func TestExit1(t *testing.T) {
	os.Exit(-1)
}

func TestExit2(t *testing.T) {
	os.Exit(-2)
}
