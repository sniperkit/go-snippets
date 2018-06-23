package main

import (
  "fmt"
  charb "github.com/gophergala/charlieswithcrayons"
)

func main() {
  fmt.Println(charb.Version)

  client := charb.New(charb.Arguments{ApiKey: "00000000-0000-0000-0000-000000000000", Url: charb.StdUrl})
  // Request for 1x random 1024bits and tag the request with 99
  bits, err := client.GetRandomBits(1024, 1, 99)
  if err != nil {
      panic(err)
  }
}
