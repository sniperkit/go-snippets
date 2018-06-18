package main

import "fmt"
import "sync"

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, num := range nums {
			out <- num
		}
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for num := range in {
			out <- num * num
		}
		close(out)
	}()
	return out
}

func merge(chans ...<-chan int) <-chan int {
  var wg sync.WaitGroup
  out := make(chan int)

  output := func(c <-chan int) {
    for n := range c {
      out<-n
    }
    wg.Done()
  }
  wg.Add(len(chans))
  for _, c := range chans {
    go output(c)
  }

  go func() {
    wg.Wait()
    close(out)
  }()
  return out
}


func main() {
	foo := gen(1, 2, 3, 4)

	bar1 := sq(foo)
	bar2 := sq(foo)

  for n := range merge(bar1, bar2) {
    fmt.Println(n)
  }

}
