package generic

import (
	"fmt"
)

type GenericError struct {
	message string
}

func (e GenericError) Error() string {
	return e.message
}

func Error(args ...interface{}) GenericError {
	e := GenericError{
		message: fmt.Sprint(args...),
	}
	return e
}

func Errorf(format string, args ...interface{}) GenericError {
	e := GenericError{
		message: fmt.Sprintf(format, args...),
	}
	return e
}
