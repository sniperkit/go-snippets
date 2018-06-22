package main

import "fmt"

type Mutable struct {
    a int
    b int
}

func (m Mutable) StayTheSame() {
    m.a = 5
    m.b = 7
}

func (m *Mutable) Mutate() {
    m.a = 5
    m.b = 7
}

func main() {
    // With method receivers that take pointers
    // Go conveniently allows both pointers and non-pointers to be passed and it automatically does the conversion
    m := &Mutable{-1, -1}
    fmt.Println(m)
    m.StayTheSame()
    fmt.Println(m)
    m.Mutate()
    fmt.Println(m)
    fmt.Println("--- using the new operator ---")
    m = new (Mutable)
    fmt.Println(m)
    m.StayTheSame()
    fmt.Println(m)
    m.Mutate()
    fmt.Println(m)
}
