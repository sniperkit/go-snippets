sign :=make(chan byte, 3)
go func(){ // G2
//...
  sign <- 2
}

go func(){ // G3
//...
  sign <- 3
}

go func(){ // G4
//...
  sign <- 4
}

for i:=0;i<3;i++{
  fmt.Printf("G%d id ended.\n", <-sign)
}

var wg sync.WaitGroup
wg.Add(3)

go func(){  // G2
  //...
  wg.Done()
}

go func(){  // G3
  //...
  wg.Done()
}

go func(){  // G4
  //...
  wg.Done()
}

wg.Wait()
fmt.Println("G2 G3 G4 are ended.")

// 我们在启用G2 G3 G4之前声明了一个sync.WaitGroup类型的变量wg，并调用其值的Add方法以使其中的计数值等于将要额外启用的Goroutine的个数。
// 后，在G2 G3 G4的运行即将结束之前，我们分别通过调用wg.Done方法将其中的计数值减去1.最后，我们在G1中调用wg.Wait方法以等待G2、G3、G4中的那3个
// wg.Done方法的调用完成。待着3个调用完成之时，在wg.Wait()处的阻塞的G1会被唤醒，它后面的那条语句也会被立即执行
// sync.WaitGroup 类型值的Add方法的第一次调用，应该发生在Done方法和Wait方法之前