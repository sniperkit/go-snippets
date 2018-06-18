package multierr

import (
	"errors"
	"testing"
)

func TestNewMultierr(t *testing.T) {
	multerr := NewMultierr()
	if multerr == nil {
		t.Error("cannot instantiate Multierr")
	}

	if multerr.Len() != 0 {
		t.Errorf("incorrect number of errors, expected %v got %v", 0, multerr.Len())
	}

	multerr.Add(errors.New("first error"))

	if multerr.Len() != 1 {
		t.Errorf("incorrect number of errors, expected %v got %v", 1, multerr.Len())
	}

	multerr.Add(errors.New("second error"))

	if multerr.Len() != 2 {
		t.Errorf("incorrect number of errors, expected %v got %v", 2, multerr.Len())
	}

	expected := "first error\nsecond error\n"
	res := multerr.Error()

	if res != expected {
		t.Errorf("incorrect error value, expected \n`%v` got \n`%v`", expected, res)
	}

	multerr = NewMultierr(errors.New("first error"), errors.New("second error"))

	if multerr.Len() != 2 {
		t.Errorf("incorrect number of errors, expected %v got %v", 2, multerr.Len())
	}
}
