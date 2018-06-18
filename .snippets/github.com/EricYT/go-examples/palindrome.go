package main

import (
  "fmt"
)

func isPalindrome(s string) bool {
  length := len(s)
  if length < 1 {
    return false
  }

  left := 0
  right := length-1

  for left < right {
    if s[left] != s[right] {
      return false
    }
    left++
    right--
  }
  return true
}

func isPalindrome2(s string) bool {
  length := len(s)
  if length < 1 {
    return false
  }

  middle := length / 2
  left := middle - 1
  right:= middle + length%2
  if left < right {
    if s[left] != s[right] {
      return false
    }
    right++
    left--
  }
  return true
}

func main() {
  var1 := "madam"
  fmt.Println(isPalindrome(var1))
  fmt.Println(isPalindrome2(var1))

  var1 = "maam"
  fmt.Println(isPalindrome(var1))
  fmt.Println(isPalindrome2(var1))

  var1 = "mcam"
  fmt.Println(isPalindrome(var1))
  fmt.Println(isPalindrome2(var1))
}
