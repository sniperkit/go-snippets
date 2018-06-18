package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
)

type Foo struct {
	name string
}

//func (f *Foo) spark() {
func (f Foo) spark() {
	fmt.Println("function call by Foo spark")
}

type Bar struct {
	name string
}

func (b Bar) spark() {
	fmt.Println("function call by Bar spark")
}

type sparkint interface {
	spark()
}

func do_spark(s interface{}) {
	sType := reflect.TypeOf(s)
	fmt.Println("do_spark type:", sType)
	fmt.Println("do_spark type name:", sType.Name())
	switch sType.Name() {
	case "Foo":
		fmt.Println("internal do_spark is main.Foo")
		//foo, ok := s.(Foo)
		foo, ok := s.(sparkint)
		if !ok {
			fmt.Println("convert foo error ", ok)
			return
		}
		foo.spark()
	case "Bar":
		fmt.Println("internal do_spark is main.Bar")
		bar, ok := s.(Bar)
		if !ok {
			fmt.Println("convert bar error ", ok)
			return
		}
		bar.spark()
	}
}

func loop(msg <-chan interface{}) {
	log.Println("Loop start ...............")
	for {
		select {
		case ops, ok := <-msg:
			if !ok {
				log.Println("loop receive error")
				continue
			}
			opsType := reflect.TypeOf(ops)
			fmt.Println("receive msg type:", opsType)
			switch opsType.Name() {
			case "Foo":
				fmt.Println("loop receive type Foo")
				foo, ok := ops.(Foo)
				if !ok {
					log.Println("convert error foo")
					continue
				}
				foo.spark()
			case "Bar":
				fmt.Println("loop receive type Bar")
				bar, ok := ops.(Bar)
				if !ok {
					log.Println("convert error bar")
					continue
				}
				bar.spark()
			default:
				log.Println("loop receive other type ", opsType.Name())
			}
		}
	}
}

func main() {
	var foo = Foo{}
	var bar = Bar{}

	fooType := reflect.TypeOf(foo)
	barType := reflect.TypeOf(bar)

	fmt.Println("foo type:", fooType.Name())
	fmt.Println("bar type:", barType.Name())

	do_spark(foo)
	do_spark(bar)

	/* chain */
	opsChan := make(chan interface{})
	go loop(opsChan)

	opsChan <- foo
	opsChan <- bar

	/* catch ctrl-c */
	interruptSig := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(interruptSig, os.Interrupt, os.Kill)

	go func() {
		s := <-interruptSig
		fmt.Println("receive signal ", s)
		done <- true
	}()

	fmt.Println("Main function end")
	<-done

}
