package main
import (
  "fmt"
  "strconv"
  "os"
)

func main() {
  // first argument starts at 0, thats why default length is 0.
  args := len(os.Args)
  if args != 2 {
    fmt.Println("exiting..." + strconv.Itoa(args))
    os.Exit(1)
  }
  fmt.Println("current directory..", os.Args[0])
}
