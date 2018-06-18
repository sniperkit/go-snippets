// Channel 是Go中的一种类型，和goroutine一起为Go提供了并发技术。
// Go鼓励人们通过Channel在goroutine之间传递数据的引用(就像把数据的owner从一个goroutine传递给另外一个goroutine), 
// Effective Go总结了这么一句话：Do not communicate by sharing memory; instead, share memory by communicating.
// 在 Go内存模型指出了channel作为并发控制的一个特性：A send on a channel happens before the corresponding receive from that channel completes. (Golang Spec)
package channel

import (
	"github.com/golang/go/src/sync"
)

// Or Channel 模式：Goroutine 方式
// 同时处理n个channel，它为每个channel启动一个goroutine，只要任意一个goroutine从channel读取到数据，输出的channel就被关闭掉了。
// 为了避免并发关闭输出channel的问题，关闭操作只执行一次。
func or(chans ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		var once sync.Once
		for _, c := range chans {
			go func(c <-chan interface{}) {
				select {
				case <-c:
					once.Do(func() { close(out) })
				case <-out:
				}
			}(c)
		}
	}()
	return out
}

// Or Channel 模式：Reflect 方式
// Go的反射库针对select语句有专门的数据(reflect.SelectCase)和函数(reflect.Select)处理。
// 所以我们可以利用反射“随机”地从一组可选的channel中接收数据，并关闭输出channel。
func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}
	orDone := make(chan interface{})
	go func() {
		defer close(orDone)
		var cases []reflect.SelectCase
		for _, c := range channels {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}
		reflect.Select(cases)
	}()
	return orDone
}

// Or Channel 模式：Recursive 方式
// 递归方式一向是比较开脑洞的实现，下面的方式就是分而治之的方式，逐步合并channel，最终返回一个channel。
func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}
	orDone := make(chan interface{})
	go func() {
		defer close(orDone)
		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			m := len(channels) / 2
			select {
			case <-or(channels[:m]...):
			case <-or(channels[m:]...):
			}
		}
	}()
	return orDone
}

// Or-Done-Channel 模式
// 这种模式是我们经常使用的一种模式，通过一个信号channel(done)来控制(取消)输入channel的处理。
// 一旦从done channel中读取到一个信号，或者done channel被关闭， 输入channel的处理则被取消。
// 这个模式提供一个简便的方法，把 done channel 和 输入 channel 融合成一个输出channel。
func orDone(done <-chan struct{}, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

// Map
// map将一个channel映射成另外一个channel，channel的类型可以不同。
func mapChan(in <-chan interface{}, fn func(interface{}) interface{}) <-chan interface{} {
	out := make(chan interface{})
	if in == nil {
		close(out)
		return out
	}
	go func() {
		defer close(out)
		for v := range in {
			out <- fn(v)
		}
	}()
	return out
}

// Reduce
// map和reduce是一组常用的操作，你可以用reduce实现`sum`、`max`、`min`等聚合操作。
func reduce(in <-chan interface{}, fn func(r, v interface{}) interface{}) interface{} {
	if in == nil {
		return nil
	}
	out := <-in
	for v := range in {
		out = fn(out, v)
	}
	return out
}

// Skip
// 集合操作，skip函数从一个channel中跳过开一些数据，然后才开始读取。
func skip(done <-chan struct{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case <-valueStream:
			}
		}
		for {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

// Take
// skip的反向操作，读取一部分数据。
func take(done <-chan struct{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}
