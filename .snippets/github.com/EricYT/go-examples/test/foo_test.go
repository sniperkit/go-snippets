package test

import (
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
)

func Test_Division_1(t *testing.T) {
	if i, e := Division(6, 2); i != 3 || e != nil {
		t.Error("not pass")
	} else {
		t.Log("pass Division")
	}
}

func TestaDiv(t *testing.T) {
	t.Log("Pass anyway")
	//t.FailNow()
	//t.Error("not pass anyway")
}

func Benchmark_Division(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Division(4, 5)
	}
}

func Benchmark_now(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pipe := make(chan int64, 1)
		test := &Test{}
		go func() {
			log.WithFields(log.Fields{
				"module": "xxx",
			}).Debugln(time.Now())
			pipe <- test.Call()
		}()
		<-pipe
	}
}
