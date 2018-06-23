package test

import (
	"syscall"
	"testing"
	"time"
)

func now() time.Time {
	var tv syscall.Timeval
	syscall.Gettimeofday(&tv)
	return time.Unix(0, syscall.TimevalToNsec(tv))
}

func BenchmarkTimeNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Now()
	}
}

func BenchmarkNowGettimeofday(b *testing.B) {
	for i := 0; i < b.N; i++ {
		now()
	}
}
