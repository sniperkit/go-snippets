/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd..
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/04/21        Liu Jiachang
 */

package task

import (
	"errors"
	"sync"
	"time"
)

// 定义错误类型
var (
	ErrHandlerExists = errors.New("This type handler is exists!")
	ErrMaxWorker     = errors.New("Max worker!")
)

// 定义处理方法类型
type Handler func(task *Task) error

// Server 配置参数
type ServerOption struct {
	queueEngine QueueEngine
	WorkerSize  int
}

// Server 结构
type Server struct {
	lock        sync.RWMutex
	mux         map[int]Handler

	distChan    chan chan Task
    resultChan  chan Result

	option      *ServerOption
}

// 创建 Server
func NewServer(option *ServerOption) *Server {
	s := Server{
		distChan:        make(chan chan Task),
		mux:             make(map[int]Handler),
		resultChan:      make(chan Result),
		option:          option,
	}

	for i := 0; i < option.WorkerSize; i++ {
		w := Worker{
			server:   &s,
			tchan:    make(chan Task),
		}

		w.Run()
	}

	return &s
}

// 注册 type 对应处理方法
func (this *Server) Register(t int, h Handler) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	_, ok := this.mux[t]

	if !ok {
		this.mux[t] = h

		return nil
	}

	return ErrHandlerExists
}

func (this *Server) FetchTasks() ([]Task, error) {
	return this.option.queueEngine.FetchTasks(this.option.WorkerSize)
}

func (this *Server) DeleteTask(id interface{}) error {
	return this.option.queueEngine.DeleteTask(id)
}

func (this *Server) Activate(id interface{}, status int16) error {
	return this.option.queueEngine.Activate(id, status)
}

// 启动 Server
func (this *Server) Start() {
	for {
		tasks, err := this.FetchTasks()

		if err != nil {
			time.Sleep(time.Second)

			tasks, err = this.FetchTasks()
		}

		for i := 0; i < len(tasks); {
			select {
			case tchan := <-this.distChan:
				this.Activate(tasks[i].Id, TaskExecuted)
				tchan <- tasks[i]
				i++
			case result := <- this.resultChan:
				if result.IsWorked == true {
					this.DeleteTask(result.Id)
				}
			}
		}
	}
}
