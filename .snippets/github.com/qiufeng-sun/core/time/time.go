package time

import (
	"time"
)

//
var (
	now      time.Time
	unix     int64
	millisec int64
)

//
func Update() {
	now = time.Now()
	unix = now.Unix()
	millisec = now.UnixNano() / 1000000
}

//
func Now() time.Time {
	return now
}

//
func Unix() int64 {
	return unix
}

//
func MillSec() int64 {
	return millisec
}
