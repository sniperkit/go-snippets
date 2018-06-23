package test

import (
	"errors"
)

func Division(a, b float64) (float64, error) {
	if 0 == b {
		return 0, errors.New("b can't be zero")
	}

	return a / b, nil
}

type Test struct {
	bar int64
}

func (t *Test) Call() int64 {
	t.bar++
	return t.bar
}
