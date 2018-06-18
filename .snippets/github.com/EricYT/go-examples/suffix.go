package main

import "fmt"

func main() {
  foo := 1
  fmt.Println("Foo is ", foo)

  foo++
  /*
  In Go,the increment statement is not a expression,and only the suffix syntax is allowed
  This line increments the variable,but it does not evaluate to anything
  */

//  ++foo

  // This is not valid
  //bar := foo++
  bar := 2

  fmt.Println("Bar is ", bar)
  fmt.Println("Foo is ", foo)
}
