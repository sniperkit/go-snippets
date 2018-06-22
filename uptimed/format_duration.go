package main

import (
	"fmt"
	"math"
	"time"
)

const (
	Second = float64(time.Second)
	Minute = float64(time.Minute)
	Hour   = float64(time.Hour)
)

func formatDuration(d time.Duration) string {
	nanos := float64(d)
	//secs := int64(math.Mod(nanos/Second, 60))
	mins := int64(math.Mod(nanos/Minute, 60))
	hours := int64(math.Mod(nanos/Hour, 24))
	return fmt.Sprintf("%.2d:%.2d", hours, mins)
}
