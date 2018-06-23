package main

//This is a test of the context package

import (
	"context"
	"fmt"
	"time"
)

//Context包为上下文管理提供了便利，
//一般情况下通过继承context.Background生成满足我们自己需求的context, 然后就可以用来管理我们的过程调用上下文了。
//这里当然不局限于web应用.

func main() {
	withValue := context.WithValue(context.Background(), "SID", "12345465")
	ctx, cancel := context.WithCancel(withValue)
	go ContextWithValueAndCancelTest(ctx, 1)
	go ContextWithValueAndCancelTest(ctx, 2)
	go ContextWithValueAndCancelTest(ctx, 3)
	go ContextWithValueAndCancelTest(ctx, 4)
	go ContextWithValueAndCancelTest(ctx, 5)
	go ContextWithValueAndCancelTest(ctx, 6)
	go ContextWithValueAndCancelTest(ctx, 7)

	<-time.Tick(time.Second * 5)
	cancel()
	<-time.Tick(time.Second * 5)

	//这里既可以通过withTimeout的超时来取消，也可以主动调用cancel回调函数来取消
	//withTimeout, cancel := context.WithTimeout(withValue, time.Duration(time.Second*13))
	withTimeout, _ := context.WithTimeout(withValue, time.Duration(time.Second*13))
	ContextWithValueAndTimeoutTest(withTimeout, 8)
	ContextWithValueAndTimeoutTest(withTimeout, 9)
	ContextWithValueAndTimeoutTest(withTimeout, 10)
	ContextWithValueAndTimeoutTest(withTimeout, 11)
	ContextWithValueAndTimeoutTest(withTimeout, 12)
	ContextWithValueAndTimeoutTest(withTimeout, 13)
	ContextWithValueAndTimeoutTest(withTimeout, 14)

	withTimeout, cancel = context.WithTimeout(withValue, time.Duration(time.Second*1300))
	go MultiLevelCancel(withTimeout, 1233445666)
	cancel()
	<-time.Tick(time.Second * 100)
}

func ContextWithValueAndCancelTest(ctx context.Context, i int32) {
	sid := ctx.Value("SID")
	sessionid, _ := sid.(string)
	fmt.Println(sessionid)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("With Cancel Cancelled", i)
			return
		}
	}
}

func ContextWithValueAndTimeoutTest(ctx context.Context, i int32) {
	sid := ctx.Value("SID")
	sessionid, _ := sid.(string)
	fmt.Println(sessionid)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("With timeout Cancelled", i)
			return
		}
	}
}

func MultiLevelCancel(ctx context.Context, i int32) {
	go func() {
		sid := ctx.Value("SID")
		sessionid, _ := sid.(string)
		fmt.Println(sessionid)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("With timeout Cancelled level 1", i)
				return
			}
		}
	}()
	MultiLevelCancelL2(ctx, i)
}

func MultiLevelCancelL2(ctx context.Context, i int32) {
	go func() {
		sid := ctx.Value("SID")
		sessionid, _ := sid.(string)
		fmt.Println(sessionid)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("With timeout Cancelled level 2", i)
				return
			}
		}
	}()
}
