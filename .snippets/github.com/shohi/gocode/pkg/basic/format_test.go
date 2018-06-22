package basic

import (
	"fmt"
	"log"
	"testing"
)

func TestFormatVsLog(t *testing.T) {
	fmt.Println("fmt first")
	log.Println("log first")
}
