package main

import "fmt"

const (
  One = iota
  Two
  Three = iota
)

const (
  Four = Three //2
  Five = iota  //1
  Six = Two    //1
  Seven = iota //3
  Eight = iota //4
)

func main() {


  fmt.Println("First ", One)
  fmt.Println("Last ", Three)

  fmt.Println("1 ", Four)
  fmt.Println("2 ", Five)
  fmt.Println("3 ", Six)
  fmt.Println("4 ", Seven)
  fmt.Println("5 ", Eight)
  //  fmt.Println(iota)    error
}
