package factory

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(fn func()) string {
	// first, we should replace stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// now the output is already logged, so invoke fn()
	fn()

	// read from the pipe and produce a string
	out := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		out <- buf.String()
	}()

	// restore stdout
	w.Close()
	os.Stdout = old

	// return captured string
	return <-out
}

func TestProfilers(t *testing.T) {
	{
		ctx := DatabaseProfilerContext{}.new("test123", "test456", 789)
		mp := DatabaseProfilerMemory{}
		mp.Post(ctx)
		mp.Post(ctx)

		if len(mp.Log) != 2 {
			t.Errorf("Expected to log 2 queries, got %d", len(mp.Log))
		}

		captureOutput(func() {
			mp.Flush()
		})

		if len(mp.Log) != 0 {
			t.Errorf("Expected empty log after flush, have %d", len(mp.Log))
		}
	}

	{
		mp := DatabaseProfilerStdout{}
		ctx := DatabaseProfilerContext{}.new("test123", "test456", 789, mp)

		output := captureOutput(func() {
			mp.Post(ctx)
		})

		if len(output) < 10 {
			t.Errorf("Unexpected length < 10 ('%s')", output)
		}

		tests := []string{" test123 ", "s]", "\"test456\"", " 789", "factory.DatabaseProfilerStdout{}"}
		for _, match := range tests {
			if !strings.Contains(output, match) {
				t.Errorf("Doesn't contain expected string %s, output '%s'", match, output)
			}
		}
	}
}
