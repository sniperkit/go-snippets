package main

import (
    "fmt"
)

type Rect struct {
    x, y float64
    width, height float64
}

func NewRect(x, y, width, height float64) *Rect {
    return &Rect{x, y, width, height}
}

func main() {
    //r1 := new(Rect)
    r2 := Rect{}
    r2.x = 15
    r3 := &Rect{0, 0, 100, 200}    
    r3.x = 10
    //r4 := &Rect{width: 200, height: 300}
    fmt.Printf("%v, %T\n", r2, r2)
    fmt.Printf("%v, %T\n", r3, r3)
    r5 := NewRect(0, 5, 100, 250)
    fmt.Printf("%v, %T\n", r5, r5)
}
