package main

// from stackover flow
// http://stackoverflow.com/questions/4220745/how-to-select-for-input-on-a-dynamic-list-of-channels-in-go

import (
	"log"
	"reflect"
)

func sendToAny(ob int, chs []chan int) int {
	set := []reflect.SelectCase{}
	for _, ch := range chs {
		set = append(set, reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(ch),
			Send: reflect.ValueOf(ob),
		})
	}
	to, _, _ := reflect.Select(set)
	return to
}

func recvFromAny(chs []chan int) (val int, from int) {
	set := []reflect.SelectCase{}
	for _, ch := range chs {
		set = append(set, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}
	from, valValue, _ := reflect.Select(set)
	val = valValue.Interface().(int)
	return
}

func main() {
	channels := []chan int{}
	for i := 0; i < 5; i++ {
		channels = append(channels, make(chan int))
	}

	go func() {
		wg.Done()
		for i := 0; i < 10; i++ {
			x := sendToAny(i, channels)
			log.Printf("Sent %v to ch%v", i, x)
		}
	}()

	for i := 0; i < 10; i++ {
		v, x := recvFromAny(channels)
		log.Printf("Received %v from ch%v", v, x)
	}
}
