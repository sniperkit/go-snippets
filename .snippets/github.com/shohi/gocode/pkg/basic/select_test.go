package basic

import (
	"fmt"
	"log"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestSelectForChannel(t *testing.T) {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			log.Println("tick.")
		case <-boom:
			log.Println("BOOM!")
			return
		default:
			log.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func TestSelectForOrder(t *testing.T) {

	// Select Not keep the order of declaration
	ch := make(chan string, 1)
	ch <- "1, not default"

	ch2 := make(chan string, 1)
	ch2 <- "2, not default"

	select {
	case v2 := <-ch2:
		log.Println(v2)
	case v1 := <-ch:
		log.Println(v1)
	default:
		log.Println("default")
	}
}

func TestSelectForDefault(t *testing.T) {
	ch := make(chan string)
	flag := false
	start := time.Now()
	log.Printf("process start at: %v\n", start)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		cnt := 0
		for {
			select {
			case <-ch:
				log.Printf("after times: %d, finally return\n", cnt)
				log.Printf("process end at: %v\n", time.Now())
				return
			default:
				if !flag {
					flag = true
				} else {
					cnt++
				}
			}
		}
	}()

	time.AfterFunc(10*time.Second, func() {
		ch <- "hello"
	})
	wg.Wait()
}

// ref, https://stackoverflow.com/questions/13666253/breaking-out-of-a-select-statement-when-all-channels-are-closed
func TestSelectWhenSenderChannelClosed(t *testing.T) {
	var ch = make(chan int)
	close(ch)

	var ch2 = make(chan int)
	go func() {
		for i := 1; i < 10; i++ {
			ch2 <- i
		}
		close(ch2)
	}()

	for i := 0; i < 20; i++ {
		select {
		case x, ok := <-ch:
			log.Println("closed", x, ok)
		case x, ok := <-ch2:
			log.Println("open", x, ok)
		}
	}
}

func TestSelectWhenRecieverChannelClosed(t *testing.T) {
	ch1 := make(chan int)
	close(ch1)
	ch2 := make(chan int)
	select {
	case ch1 <- 10:
		log.Println("ch1 gets value")
	case ch2 <- 20:
		log.Println("ch2 gets value")
	}
}

func TestSelectWithReflect(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)
	chanMap := map[interface{}]struct{}{
		ch1: struct{}{},
		ch2: struct{}{},
	}

	go func() {
		ch1 <- 10
	}()

	go func() {
		ch2 <- 20
	}()

	defer func() {
		var cases []reflect.SelectCase
		for c := range chanMap {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c)})
		}
		// We exepct one read from each block
		for i := 0; i < len(cases); i++ {
			reflect.Select(cases)
		}
	}()
}

// ref,
// https://play.golang.org/p/1ZVdYlA52x
func TestPriSelect(t *testing.T) {

	// Helper functions
	generate := func(c chan int, v int) {
		for {
			c <- v
		}
	}

	priselect := func(a, b, c, d chan int) int {
		// priselect {
		// case x := <-a:
		// 	return x
		// case x := <-b:
		// 	return x
		// case x := <-c:
		// 	return x
		// case x := <-d:
		// 	return x
		// }

		select {
		case x := <-a:
			return x
		default:
		}
		select {
		case x := <-a:
			return x
		case x := <-b:
			return x
		default:
		}
		select {
		case x := <-a:
			return x
		case x := <-b:
			return x
		case x := <-c:
			return x
		default:
		}
		select {
		case x := <-a:
			return x
		case x := <-b:
			return x
		case x := <-c:
			return x
		case x := <-d:
			return x
		}
		panic("unreachable")
	}

	// Test codes
	a, b, c, d := make(chan int), make(chan int), make(chan int), make(chan int)
	go generate(a, 0)
	go generate(b, 1)
	go generate(c, 2)
	go generate(d, 3)
	count := make([]int, 4)
	for i := 0; i < 10000; i++ {
		count[priselect(a, b, c, d)]++
	}
	fmt.Printf("%v\n", count)
}
