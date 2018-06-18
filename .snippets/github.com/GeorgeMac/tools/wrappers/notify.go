package wrappers

import (
	"os"
	"os/signal"
)

func Notify(c chan<- struct{}, sig ...os.Signal) {
	go func() {
		ch := make(chan os.Signal, cap(c))
		signal.Notify(ch, sig...)
		for range ch {
			c <- struct{}{}
		}
	}()
}
