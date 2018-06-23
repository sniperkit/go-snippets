/**
 * 无缓冲通道：使用通道，在 goroutine 之间进行数据同步
 * 使用2个 goroutine 来模拟网球比赛，并使用无缓冲的通道来模拟球的来回
 */
package main

import (
	"sync"
	"math/rand"
	"time"
	"fmt"
)

var wait sync.WaitGroup

func init() {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// 创建无缓冲通道
	court := make(chan int)

	wait.Add(2)

	// 启动2个球手
	go player("James", court)
	go player("Lily", court)

	// 发球
	court <- 1

	// 等待游戏结束
	wait.Wait()
}

// player 模拟一个选手在打网球
func player(name string, court chan int) {
	defer wait.Done()

	for {
		ball, ok := <-court
		if !ok {
			// 如果通道关闭，我们就赢了
			fmt.Printf("Player %s Won\n", name)
			return
		}

		// 选随机数，然后用这个数来判断我们是否丢球
		n := rand.Intn(100)
		if n%13 == 0 {
			fmt.Printf("Player %s Missed\n", name)

			// 关闭通道，我们输了
			close(court)
			return
		}

		// 显示击球数，并将击球数加1
		fmt.Printf("Player %s Hit %d\n", name, ball)
		ball++

		// 将球打向对手
		court <- ball
	}
}
