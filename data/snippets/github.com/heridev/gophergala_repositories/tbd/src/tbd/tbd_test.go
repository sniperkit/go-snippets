package main

import (
	"strings"
	"testing"
)

var (
	tbdata string = `implement more #tests
extract code out of #main
support cgi/fast-cgi in #main`
)

func TestNoArgs(t *testing.T) {
	// Configure the handlers ...
	// Use a collecting handler to catch all tasks right after they are parsed
	// so we have something to easily compare against the final stage's output
	collectTest, outputTest := collecting()
	collectFinal, outputFinal := collecting()
	handlers := []handler{collectTest, counting(), tracing(), cutoff(), collectFinal}

	// ... and execute them against the sample tbdata
	reader := strings.NewReader(tbdata)
	exit := process(reader, handlers)

	if exit != 0 {
		t.Error("Unexpected exit code.", exit)
	}

	initial := outputTest()
	result := outputFinal()
	if len(result) != 2 {
		t.Error("Unexpected result length", 2)
	}
	if result[0] != initial[0] {
		t.Error("Unexpected task at result[0]", result[0])
	}
	if result[1] != initial[1] {
		t.Error("Unexpected task at result[1]", result[1])
	}
}

func TestArgs(t *testing.T) {
	// Configure the handlers ...
	// Use a collecting handler to catch all tasks right after they are parsed
	// so we have something to easily compare against the final stage's output
	collectTest, outputTest := collecting()
	collectFinal, outputFinal := collecting()
	handlers := []handler{collectTest, counting(), tracing(), matching([]string{"main"}), collectFinal}

	// ... and execute them against the sample tbdata
	reader := strings.NewReader(tbdata)
	exit := process(reader, handlers)

	if exit != 0 {
		t.Error("Unexpected exit code.", exit)
	}

	initial := outputTest()
	result := outputFinal()
	if len(result) != 2 {
		t.Error("Unexpected result length", 2)
	}
	if result[0] != initial[1] {
		t.Error("Unexpected task at result[0]", result[0])
	}
	if result[1] != initial[2] {
		t.Error("Unexpected task at result[1]", result[1])
	}
}
