package main

import "fmt"

type Foo struct {
  Bar *[]string
}

func main() {
  foo := &Foo{}

  str := []string{"abc", "def", "aaa"}

  foo.Bar = &str

  fmt.Println("---> Foo str:", *foo.Bar)
}
