package main
import "fmt"
import (
  "strconv"
)

func vals() (int, int, int) {
    return 3, 7, 10
}

func add(x int, y int) int {
    return x + y
}

func main() {
 // This means ignore the first param and second params
  _, _, c := vals()
  fmt.Println("Function value:" + strconv.Itoa(c))

  var x int
  fmt.Println("(enter x) : ")
  _, err := fmt.Scanf("%d", &x)
  if err != nil  {
        fmt.Println("Error : ", err)
  }

  var y int
  fmt.Println("(enter y) : ")
  _, err2 := fmt.Scanf("%d", &y)
  if err2 != nil  {
        fmt.Println("Error : ", err2)
  }
   sum := add(x, y)
  fmt.Println(strconv.Itoa(x) + " + " + strconv.Itoa(y) + " = " + strconv.Itoa(sum))

}
