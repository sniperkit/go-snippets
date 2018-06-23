package main

import (
  "fmt"
  "github.com/koding/kite"
)

func main() {
  k := kite.New("exp2", "1.0.0")
  // Connect to our math kite
  mathWorker := k.NewClient("http://localhost:8080/kite")
  mathWorker.Dial()
  response, _ := mathWorker.Tell("square", 4) // call "square" method with argument 4
  fmt.Println("result:", response.MustFloat64())
}
