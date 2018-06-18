package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

type MyFile struct {
	f *os.File
}

func NewFile(name string) (*MyFile, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	runtime.SetFinalizer(f, func(f *os.File) {
		fmt.Println("New file is auto close")
		f.Close()
	})
	return &MyFile{f: f}, nil
}

func test() {
	f, _ := NewFile("file.go")
	fmt.Printf("f is:%+v\n", f)
	f.f.Close()
}

func main() {
	go func() {
		test()
	}()

	time.Sleep(time.Second * 1)
	runtime.GC()

	time.Sleep(time.Second * 10)
}
