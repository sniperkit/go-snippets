// Go maneja varios tipos de valores incluyendo cadenas,
// enteros, flotantes, booleanos, etc. Aquí algunos
// ejemplos básicos.

package main

import "fmt"

func main() {

    // Cadenas, que pueden ser concatenadas con `+`.
    fmt.Println("go" + "lang")

    // Enteros y flotantes
    fmt.Println("1+1 =", 1+1)
    fmt.Println("7.0/3.0 =", 7.0/3.0)

    // Booleanos, con operadores booleanos, tal y como lo esperarías.
    fmt.Println(true && false)
    fmt.Println(true || false)
    fmt.Println(!true)
}
