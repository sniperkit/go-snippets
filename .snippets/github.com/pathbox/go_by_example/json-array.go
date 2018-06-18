package main

import(
  "fmt"
  "encoding/json"
)

func main() {
  type T struct {
    Endpoint string `json:"endpoint"`
    Counter  string `json:"counter"`
    Dstype   string `json:"dstype"`
    Step     int    `json:"step"`
    Value    []struct {
      Timestamp  int  `json:"timestamp"`
      Value      float64  `json:"value"`
    } `json:"Values"`
  }
  t := new(T)
  err := json.Unmarshal([]byte(s), t)
  if err != nil {
    fmt.Println(err)
  }
  fmt.Printf("%+v", *t)
}
