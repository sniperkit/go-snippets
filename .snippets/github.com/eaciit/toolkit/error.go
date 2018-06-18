package toolkit

import (
	"fmt"
)

func Error(txt string) error {
	return fmt.Errorf(txt)
}

func Errorf(txt string, obj ...interface{}) error {
	return fmt.Errorf(txt, obj...)
}
