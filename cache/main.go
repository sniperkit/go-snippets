package main

import (
  "fmt"
  "cache"
  "time"
  "os"
//  "path"
)

func main() {
  fmt.Println("---------------> main 1")

  cwd, err := os.Getwd()
  if err != nil {
    fmt.Println("-------------> main 2:", cwd)
  }
  err = cache.Set("test1", "Hello world", time.Second*300)
  if err != nil {
    fmt.Println("--------------> main 3:", err)
  }

  var res string
  cache.Get("test1", &res)
  fmt.Println("res:", res)
}
