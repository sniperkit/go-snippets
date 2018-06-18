package main

import "fmt"

type long int64

func (l long) Show() long { return l }

func (l long) Increment() { l++ }

type Long interface {
  Show() long
  Increment()
}

func main() {
  var l long = 231

  fmt.Println("l is ", l.(Long).Show())

  var lIf Long = l
  lIf.Increment()

  fmt.Println("l is ", lIf.Show())

}
