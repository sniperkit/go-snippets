package main

import "fmt"
import "reflect"

func main() {
  var x float64 = 3.4
  xType := reflect.TypeOf(x)
  xValue := reflect.ValueOf(x)

  fmt.Println("x type of:", xType)
  fmt.Println("x value of:", xValue)

  fmt.Println("x value type of:", xValue.Type())
  fmt.Println("x value type of:", xValue.Kind())
  fmt.Println("x value value of:", xValue.Float())

  var y int = 3234
  p := reflect.ValueOf(&y)
  fmt.Println("y type of:", p.Type())
  fmt.Println("y setability of:", p.CanSet())

  v := p.Elem()
  fmt.Println("v setability of:", v.CanSet())

  v.SetInt(23324)
  fmt.Println(v.Interface())
  fmt.Println(y)

}
