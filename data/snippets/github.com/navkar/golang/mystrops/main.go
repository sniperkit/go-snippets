package mystrops

import (
  "fmt"
)

func Reverse(s string) (string) {
  fmt.Println("input : " + s)
  result := ""
  for _,v := range s {
      result = string(v) + result
    }
  return result
}
