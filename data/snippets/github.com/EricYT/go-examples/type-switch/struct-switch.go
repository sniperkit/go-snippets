package main

import (
	"log"
	"time"
)

type emptyCtx int

func (*emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*emptyCtx) Done() <-chan struct{} {
	return nil
}

func (*emptyCtx) Err() error {
	return nil
}

func (*emptyCtx) Value(key interface{}) interface{} {
	return nil
}

func (e *emptyCtx) String() string {
	// so special switch case
	switch e {
	case background:
		return "context.Background"
	case todo:
		return "context.todo"
	}
	return "unknow empty context"
}

var (
	background = new(emptyCtx)
	todo       = new(emptyCtx)
)

func main() {
	log.SetFlags(0)
	log.Println("struct type switch go")

	log.Println(background.String())
	log.Println(todo.String())

	// foo
	var foo *emptyCtx = new(emptyCtx)
	log.Println(foo.String())
}
