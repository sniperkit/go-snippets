package bamboo

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func Foo(ctx context.Context) (context.Context, string) {
	log.Println("bamboo_test: foo function")
	return context.WithValue(ctx, "foo", 1), "foo-key"
}

func Bar(ctx context.Context, index int, key string) context.Context {
	log.Println("bamboo_test: bar run key: ", key, " index: ", index)
	return ctx
}

func TestBamboo(t *testing.T) {
	b := NewBamboo()
	err := b.Join(Foo).Join(Bar, 123).Go()
	if err != nil {
		t.Fatalf("bamboo_test: function go error: %s", err)
	}
	foo := Foo
	bar := Bar
	b1 := NewBamboo()
	err = b1.Join(&foo).Join(&bar, 123).Go()
	if err != nil {
		t.Fatalf("bamboo_test: function ptr go error: %s", err)
	}
}

func TestBambooJoin(t *testing.T) {
	b := NewBamboo()
	err := b.Join(Foo).Join(123).Go()
	switch err {
	case ErrorBambooJoinNotFunc:
	default:
		t.Fatalf("bamboo_test: join 123 should return error ErrorBambooJoinNotFunc")
	}
}

func TestBambooJoinParamsIn(t *testing.T) {
	b := NewBamboo()
	err := b.Join(Foo).Join(Bar).Go()
	switch err {
	case ErrorBambooJoinFuncParamsInNotMatch:
	default:
		t.Fatalf("bamboo_test: join Bar params not enough")
	}
	b1 := NewBamboo()
	err = b1.Join(func(index int) context.Context { fmt.Println("input params not contains context.Context"); return nil }).Go()
	switch err {
	case ErrorBambooJoinFuncFirstInParamNotContext:
	default:
		t.Fatalf("bamboo_test: join function first param is not context.Context")
	}
}

func TestBambooJoinParamsOut(t *testing.T) {
	b := NewBamboo()
	err := b.Join(Foo).Join(func(ctx context.Context) { fmt.Println("result not contains context.Context") }).Go()
	switch err {
	case ErrorBambooJoinFuncResultEmpty:
	default:
		t.Fatalf("bamboo_test: join lambda function not contains context.Context")
	}
}

func TestBambooCancel(t *testing.T) {
	b := NewBamboo()
	f1 := func(ctx context.Context) context.Context {
		time.Sleep(time.Second * 3)
		return ctx
	}
	var errCh chan error = make(chan error, 1)
	go func() {
		errCh <- b.Join(f1).Join(f1).Go()
	}()

	b.Cancel()

	select {
	case err := <-errCh:
		switch err {
		case ErrorBambooCancel:
		default:
			t.Fatalf("bamboo_test: operation already canceled.")
		}
	}
}
