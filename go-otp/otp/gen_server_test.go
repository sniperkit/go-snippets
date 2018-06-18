package otp

import (
	"context"
	"log"
	"testing"
	"time"
)

type fake struct {
}

func (f *fake) Init(args ...interface{}) (context.Context, error) {
	log.Printf("fake initialize args: %#v\n", args)
	return context.TODO(), nil
}

func (f *fake) HandleCall(ctx context.Context, req interface{}) (context.Context, interface{}, error) {
	log.Printf("fake handle call %#v\n", req)
	ctx = context.WithValue(ctx, "call", 1)
	return ctx, "pong", nil
}

func (f *fake) HandleInfo(ctx context.Context, req interface{}) (context.Context, error) {
	log.Printf("fake handle info %#v\n", req)
	ctx = context.WithValue(ctx, "info", 2)
	return ctx, nil
}

func (f *fake) HandleCast(ctx context.Context, req interface{}) (context.Context, error) {
	log.Printf("fake handle cast %#v\n", req)
	ctx = context.WithValue(ctx, "cast", 3)
	return ctx, nil
}

func (f *fake) Terminate(ctx context.Context, reason interface{}) error {
	log.Printf("fake terminate reason: %#v\n", reason)
	return nil
}

func TestGenServerCall(t *testing.T) {
	fakeName := Name("fake")
	f := &fake{}
	if err := NewGenServer(fakeName, f, 1, 2, 3); err != nil {
		t.Errorf("test gen server call new gen server error: %s\n", err)
	}
	defer Terminate(fakeName, "go die")

	res, err := Call(fakeName, "ping")
	if err != nil {
		t.Errorf("test gen server call error: %s\n", err)
	}
	log.Printf("test gen server call return: %#v\n", res)
}

func TestGenServerCast(t *testing.T) {
	fakeName := Name("fake")
	f := &fake{}
	if err := NewGenServer(fakeName, f, 1, 2, 3); err != nil {
		t.Errorf("test gen server cast new gen server error: %s\n", err)
	}
	defer Terminate(fakeName, "go die")

	if err := Cast(fakeName, "cast go"); err != nil {
		t.Errorf("test gen server cast error: %s\n", err)
	}
	time.Sleep(1 * time.Second)
}

func TestGenServerInfo(t *testing.T) {
	fakeName := Name("fake")
	f := &fake{}
	if err := NewGenServer(fakeName, f, 1, 2, 3); err != nil {
		t.Errorf("test gen server info new gen server error: %s\n", err)
	}
	defer Terminate(fakeName, "go die")

	if err := Info(fakeName, "info go"); err != nil {
		t.Errorf("test gen server info error: %s\n", err)
	}
	time.Sleep(1 * time.Second)
}
