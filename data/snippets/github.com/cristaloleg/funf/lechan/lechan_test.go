package lechan_test

import (
	"testing"
	"time"

	"github.com/cristaloleg/funf/lechan"
)

func TestLechan_New(t *testing.T) {
	ch := lechan.New(10)
	if ch == nil {
		t.Fatal("cannot instantiate Lechan")
	}
}

func BenchmarkLechan_Put(b *testing.B) {
	b.ReportAllocs()
	ch := lechan.New(1000)

	go devnull(ch)

	for i := 0; i < b.N; i++ {
		ch.Put(0)
	}
}

func BenchmarkLechan_TryPut(b *testing.B) {
	b.ReportAllocs()
	ch := lechan.New(1000)

	go devnull(ch)

	for i := 0; i < b.N; i++ {
		ch.TryPut(0)
	}
}

func BenchmarkLechan_TryCountedPut(b *testing.B) {
	b.ReportAllocs()
	ch := lechan.New(1000)

	go devnull(ch)

	for i := 0; i < b.N; i++ {
		ch.TryCountedPut(0, 3)
	}
}

func BenchmarkLechan_TryTimedPut(b *testing.B) {
	b.ReportAllocs()
	ch := lechan.New(1000)

	go devnull(ch)

	for i := 0; i < b.N; i++ {
		ch.TryTimedPut(0, 10*time.Nanosecond)
	}
}

func BenchmarkLechan_TryLimitedPut(b *testing.B) {
	b.ReportAllocs()
	ch := lechan.New(1000)

	go devnull(ch)

	for i := 0; i < b.N; i++ {
		ch.TryLimitedPut(0, 3, 10*time.Nanosecond)
	}
}

func devnull(ch *lechan.Lechan) {
	_, _ = ch.Get()
}
