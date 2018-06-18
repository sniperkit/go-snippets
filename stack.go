package main

import "fmt"

type stackEntry struct {
  next *stackEntry
  value interface{}
}

type Stack struct {
  top *stackEntry
}

func (s *Stack) Push(value interface{}) {
  var e stackEntry
  e.value = value
  e.next = s.top
  s.top = &e
}

func (s *Stack) Pop() interface{} {
  if s.top == nil {
    return nil
  }

  v := s.top.value
  s.top = s.top.next
  return v
}

func main() {
//  stack := &Stack{}
  stack := new(Stack)

  stack.Push(1)
  stack.Push(2)
  stack.Push(3)
  stack.Push("hello")
  stack.Push("world")

  fmt.Println("Stack :", stack)

  var total int
  var strTotal string
  for {
    value := stack.Pop()
    if value != nil {
      fmt.Println("Stack pop:", value)
      switch value.(type) {
      case int:
        fmt.Println("Int:", value)
        total += value.(int)
      case string:
        fmt.Println("String:", value)
        strTotal += " " + value.(string)
      }
      continue
    }
    break
  }
  fmt.Println("total:", total)
  fmt.Println("strTotal:", strTotal)
}

