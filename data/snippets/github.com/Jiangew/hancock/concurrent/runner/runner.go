// runner 包使用通道管理处理任务的运行和生命周期
package runner

import (
	"os"
	"time"
	"errors"
	"os/signal"
)

// Runner 在给定的超时时间内执行一组任务，并且在操作系统发送中断信号时结束这些任务
type Runner struct {
	// 报告从操作系统发送的信号
	interrupt chan os.Signal

	// 报告处理任务已经完成
	complete chan error

	// 报告处理任务已经超时
	timeout <-chan time.Time

	// tasks 持有一组以索引顺序依次执行的函数
	tasks []func(int)
}

// ErrTimeout 会在任务执行超时时返回
var ErrTimeout = errors.New("received timeout")

// ErrInterrupt 会在接收到操作系统的事件时返回
var ErrInterrupt = errors.New("received interrupt")

// New 返回一个新的准备使用的 Runner
func New(d time.Duration) *Runner {
	return &Runner{
		interrupt: make(chan os.Signal, 1),
		complete:  make(chan error),
		timeout:   time.After(d),
	}
}

// Add 将一个任务附加到 Runner 上；这个任务是一个接收一个 int 类型的 id 作为参数的函数
func (r *Runner) Add(tasks ...func(int)) {
	r.tasks = append(r.tasks, tasks...)
}

// Start 用不同的 goroutine 处理任务
func (r *Runner) Start() error {
	// 接收所有中断信号
	signal.Notify(r.interrupt, os.Interrupt)

	// 用不同的 goroutine 执行不同的任务
	go func() {
		r.complete <- r.run()
	}()

	// 当任务处理「完成 && 程序运行超时」时发出的信号
	select {
	case err := <-r.complete:
		return err

	case <-r.timeout:
		return ErrTimeout
	}
}

// getInterrupt 验证是否接收到了中断信号
func (r *Runner) getInterrupt() bool {
	select {
	// 当中断事件被触发时发出的信号
	case <-r.interrupt:
		// 停止接收后续的任何信号
		signal.Stop(r.interrupt)
		return true

	default:
		return false
	}
}

// run 执行每一个已注册的任务
func (r *Runner) run() error {
	for id, task := range r.tasks {
		// 检测操作系统的中断信号
		if r.getInterrupt() {
			return ErrInterrupt
		}

		// 执行注册的任务
		task(id)
	}

	return nil
}
