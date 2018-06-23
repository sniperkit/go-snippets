//  求平方数
// 第一阶段是 gen 函数，它能够将一组整数转换为 channel ，
// channel 可以将数字发送出去。 gen 函数首先启动一个 goroutine ，
// 该 goroutine 发送数字到 channel ，当数字发送完时关闭 channel
func gen(nums ...int) <-chan int{
  out := make(chan int)
  go func ()  {
    for _, n := range nums {
      out <-n
    }
    close(out)  //当数字发送完时关闭 channel
  }()
  return out
}

 // sq 函数，它从 channel 接收一个整数，然后返回 一个 channel ，返回的 channel
 // 可以发送 接收到整数的平方。 当它的 inbound channel 关闭，并且把所有数字均发送到下游时，会关闭 outbound channel

 func sq(in <-chan int) <-chan int{
   out := make(chan int)
   go func() {
     for n:= range in{
       out <- n*n
     }
     close(out) // 当数字发送完时关闭 channel
   }()
   return out
 }

 //main 函数 用于设置流水线并运行最后一个阶段。最后一个阶段会从第二阶段接收数字，并逐个打印出来，
 // 直到来自于上游的 inbound channel 关闭

 func main() {
   // 设置流水线
   c := gen(2, 3)
   out := sq(c)

   // 消费输出结果
   fmt.Println(<-out) // 4
   fmt.Println(<-out) // 9
}

// 由于 sq 函数的 inbound channel 和 outbound channel 类型一样，
// 所以组合任意个 sq 函数。比如像下面这样使用：

func main() {
    // 设置流水线并消费输出结果
    for n := range sq(sq(gen(2, 3))) {
        fmt.Println(n) // 16 then 81
    }
}


// 扇出：同一个 channel 可以被多个函数读取数据，直到 channel 关闭。 这种机制允许将工作负载分发到一组 worker ，以便更好地并行使用 CPU 和 I/O 。

// 扇入：多个 channel 的数据可以被同一个函数读取和处理，然后合并到一个 channel ，直到所有 channel 都关闭。

// 我们修改一下上个例子中的流水线，这里我们运行两个 sq 实例，它们从同一个 channel 读取数据。 这里我们引入一个新函数 merge 对结果进行"扇入"操作
func main() {
  in := gen(2,3)

  //启动两个 sq 实例，即两个 goroutines 处理 channel "in" 的数据
  c1 := sq(in)
  c2 := sq(in)
  // merge 函数将 channel c1 和 c2 合并到一起，这段代码会消费 merge 的结果

  for n:= range merge(c1, c2){
    fmt.Println(n) // 打印　４，９　或　９，４
  }
}

// merge 函数 将多个 channel 转换为一个 channel ，它为每一个 inbound channel 启动一个 goroutine ，
// 用于将数据 拷贝到 outbound channel 。 merge 函数的实现见下面代码 (注意 wg 变量)：

func merge(cs ...<-chan int) <-chan int {
  var wg sync.WaitGroup // 同步锁计数的实现变量
  out := make(chan int)

  // 为每一个输入 channel cs 创建一个 goroutine output
  // output 将数据从 c 拷贝到 out ，直到 c 关闭，然后 调用 wg.Done

  output := func(c <-chan int){
    for n:= range c {
      out <- n
    }
    wg.Done()
  }
  wg.Add(len(cs))
  for _, c := range cs {
    go output(c)
  }
  // 启动一个 goroutine ，用于所有 output goroutine 结束时，关闭 out
   // 该 goroutine 必须在 wg.Add 之后启动
   go func() {
     wg.Wait()
     close(out)
   }()
   return out
}


// 如果在创建 channel 时就知道要发送的值的个数，使用 buffer 就能够简化代码。 仍然使用求平方数的例子，
// 我们对 gen 函数进行重写。我们将这组整型数拷贝到一个 缓冲 channel 中，从而避免创建一个新的 goroutine
func gen(nums ...int) <-chan int {
    out := make(chan int, len(nums))
    for _, n := range nums {
        out <- n
    }
    close(out)
    return out
}

// 回到 流水线中被阻塞的 goroutine ，我们考虑让 merge 函数返回一个缓冲管道：

func merge(cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int, 1) // 在本例中存储未读的数据足够了
    // ... 其他部分代码不变 ...

    // 尽管这种方法解决了这个程序中阻塞 goroutine 的问题，但是从长远来看，它并不是好办法。 缓存大小选择为 1 是建立在两个前提之上：
    //
    // 我们已经知道 merge 函数有两个 inbound channel
    // 我们已经知道下游阶段会消耗多少个值
    // 这段代码很脆弱。如果我们在传入一个值给 gen 函数，或者下游阶段读取的值变少， goroutine 会再次被阻塞。
    //
    // 为了从根本上解决这个问题，我们需要提供一种机制，让下游阶段能够告知上游发送者停止接收的消息。 下面我们看下这种机制。


    // 显式取消 (Explicit cancellation)
    //
    // 当 main 函数决定退出，并停止接收 out 发送的任何数据时，它必须告诉上游阶段的 goroutine 让它们放弃 正在发送的数据。
    // main 函数通过发送数据到一个名为 done 的 channel 实现这样的机制。 由于有两个潜在的 发送者被阻塞，它发送两个值

func main() {
  in := gen(2,3)

  // 启动两个运行 sq 的 goroutine
  // 两个 goroutine 的数据均来自于 in

  c1 := sq(in)
  c2 := sq(in)

  // 消耗 output 生产的第一个值
  done := make(chan struct{}, 2)
  out := merge(done, c1, c2)
  fmt.Println(<-out) // 4 or 9

  // 告诉其他发送者，我们将要离开
  // 不再接收它们的数据
  done <- struct{}{}
  done <- struct{}{}
}

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    // 为 cs 的的每一个 输入 channel
    // 创建一个 goroutine 。 output 函数将
    // 数据从 c 拷贝到 out ，直到 c 关闭，
    // 或者接收到 done 信号；
    // 然后调用 wg.Done()
    output := func(c <-chan int){
      for n := range c {
        select{
        case out <-n:
        case <-done:
        }
        wg.Done()
      }
      // the rest is unchanged
    }

}

// 我们只要关闭 done channel ，就能够让解开对所有发送者的阻塞。
// 对一个channel的关闭操作事实上是对所有接收者的广播信号

// 我们把 done channel 作为一个参数传递给每一个 流水线上的函数，通过 defer 表达式声明对 done channel 的关闭操作。
// 因此，所有从 main 函数作为源头被调用的函数均能够收到 done 的信号，每个阶段都能够正常退出。 使用 done 对 main 函数重构以后

func main() {
    // 设置一个 全局共享的 done channel ，
    // 当流水线退出时，关闭 done channel
    // 所有 goroutine 接收到 done 的信号后，
    // 都会正常退出。
    done := make(chan struct{})
    defer close(done)

    in := gen(done, 2, 3)

    // 将 sq 的工作分发给两个 goroutine
    // 这两个 goroutine 均从 in 读取数据
    c1 := sq(done, in)
    c2 := sq(done, in)

    // 消费 outtput 生产的第一个值
    out := merge(done, c1, c2)
    fmt.Println(<-out) // 4 or 9

    // defer 调用时， done channel 会被关闭。
}

// 现在，流水线中的每个阶段都能够在 done channel 被关闭时返回。 merge 函数中的 output 代码也能够顺利返回，因为它 知道 done channel 关闭时，
// 上游发送者 sq 会停止发送数据。 在 defer 表达式执行结束时，所有调用链上的 output 都能保证 wg.Done() 被调用

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    // 为 cs 的每一个 channel 创建一个 goroutine
    // 这个 goroutine 运行 output ，它将数据从 c
    // 拷贝到 out ，直到 c 关闭，或者 接收到 done
    // 的关闭信号。人啊后调用 wg.Done()
    output := func(c <-chan int) {
        defer wg.Done()
        for n := range c {
            select {
            case out <- n:
            case <-done:
                return
            }
        }
    }
    // ... the rest is unchanged ...
}

同样的原理， done channel 被关闭时， sq 也能够立即返回。在 defer 表达式执行结束时，所有调用链上的 sq 都能保证 out channel 被关闭


func sq(done <-chan struct{}, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            select {
            case out <- n * n:
            case <-done:
                return
            }
        }
    }()
    return out
} 
