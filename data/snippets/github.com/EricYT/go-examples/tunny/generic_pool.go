package main

import (
	"log"
	"sync"
	"time"

	"github.com/Jeffail/tunny"
)

type Foo struct {
	key []byte
	val int
}

type Bar struct {
	foo *Foo
}

func (b *Bar) Run() {
	//tmp := (*b).foo.val + 1
	//(*b).foo.val = tmp
	tmp := (*b).foo.val + 1
	b.foo.val = tmp

	key := "bar"
	b.foo.key = []byte(key)
}

func NewFoo(i int, res *Foo) func() {
	b := &Bar{
		foo: res,
	}
	return b.Run
}

func main() {
	pool, err := tunny.CreatePoolGeneric(10).Open()
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	var wg sync.WaitGroup
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go func(idx int) {
			defer wg.Done()
			var result = Foo{val: idx}
			_, err := pool.SendWorkTimed(time.Second, NewFoo(idx, &result))
			if err != nil {
				panic(err)
			}
			// result
			log.Printf("idx: %d key: %s reuslt: %d\n", idx, string(result.key), result.val)
			if result.val != (idx + 1) {
				panic("result should be added one")
			}
		}(i)
	}

	wg.Wait()
}
