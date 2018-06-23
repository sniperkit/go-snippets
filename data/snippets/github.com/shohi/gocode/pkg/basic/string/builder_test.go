package string

import (
	"fmt"
	"strings"
	"testing"
)

func TestStringBuilder(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("Hello")
	sb.WriteByte(byte(' '))
	sb.WriteString("World")

	fmt.Printf("content: %s\n", sb.String())
}
