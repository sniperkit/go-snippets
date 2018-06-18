package lib

import (
	"fmt"
)

type S1 struct{}

type Interfacer interface {
	read(b string, i, j, k int) string
}

func (s *S1) read(b string, i, j, k int) string {

}
