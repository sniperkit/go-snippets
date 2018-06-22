package main
import
( "fmt"
  "strconv")

func main() {
  // using fmt.Scanf() with format %q  a double-quoted
  // string safely escaped with Go syntax
  var input string
  fmt.Println("--- Demo of Scanf with %q ---")
  fmt.Println("(enter a quoted string) : ")
  // The scanf requires a new line character to work for the next scanf
  fmt.Scanf("%q\n", &input)
  fmt.Println("[input entered : " + input + " ]")
  var amount int
  fmt.Println("(enter amount) : ")
  _, err := fmt.Scanf("%d", &amount)
  if err != nil  {
        fmt.Println("Error : ", err)
  }
  fmt.Println("[amount entered : " + strconv.Itoa(amount) + " ]")
  return
}
