package time

import (
	"log"
	"strings"
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	d := 1000 * time.Second
	log.Println(strings.ToUpper(d.String()))
}

func TestParseDuration(t *testing.T) {
	durationStr := "10s"
	log.Println(time.ParseDuration(durationStr))

	durationStr = "10"
	log.Println(time.ParseDuration(durationStr))
}

func TestDurationZeroValue(t *testing.T) {
	var d time.Duration
	log.Printf("zero value of duration is %v", d)
}

func TestDurationFromFloat64(t *testing.T) {
	//
	aa := 17.58
	log.Printf("duration: %v", int64(aa))
	log.Printf("duration: %v", time.Duration(aa)*time.Second)
}
