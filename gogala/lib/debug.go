package lib

import "log"

type Debug bool

func (d Debug) Printf(format string, args ...interface{}) {
	if d {
		log.Printf(format, args...)
	}
}
