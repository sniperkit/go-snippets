package main

import (
	stderrs "errors"

	"github.com/albert-widi/go_common/errors"
	"github.com/albert-widi/go_common/logger/go-kit"
)

func main() {
	l := logger.New()
	errors.SetRuntimeOutput(true)
	err := errors.New("This is an error", errors.Fields{"field1": "value1"})
	err2 := stderrs.New("This is an error")
	if errors.Match(err, err2) {
		l.Errors(err)
	}
}
