/******************************************
 *                                        *
 *   Is interface{} a generic pointer ?   *
 *                                        *
 ******************************************/

package main

import "fmt"

func main() {
    var p interface{}

    // As a string
    p = "Hello, Planet!"
    fmt.Println(p)

    // As a pointer to a string
    str := "Oops!"
    p = &str
    fmt.Println(p, *(p.(*string)))

    // As a int
    p = 123
    fmt.Println(p)

    // As a Imaginary
    p = 1 + 1i
    fmt.Println(p)

    // As a pointer to a function
    p = func(x float64, y float64) float64 {
        return (x + y) / 2.0
    }
    a, b := 7.0, 9.0
    fmt.Println("average of ", a, ", ", b, " is ", p.(func(float64, float64) float64)(a, b))

    // As a pointer to a variadic function
    p = func(args ...int) int {
        total := 0
        for _, v := range args {
            total += v
        }
        return total
    }
    slice := []int{1, 2, 3, 4, 5}
    fmt.Println("Sum of ", slice, " is ", p.(func(...int) int)(slice...))
}
