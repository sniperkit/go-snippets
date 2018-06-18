var c = make(chan bool)
func Actor() {
  <-c  // 观察者
  // doing something
}

func main() {
  go Actor() // 观察者在后台开始准备
  c <- true  // 通知观察者
}

// 一个goroutine就是一个actor，通信是通过语言提供的管道完成的。go语言的goroutine是非常轻量级的，又可以充分发挥多核的优势。
// actor模式的核心就在这里，无锁+充分利用多核，actor之间通过消息通信共同工作。