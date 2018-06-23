// `for` es la única estructura de control iterativa en Go.
// Hay tres tipos básicos de ciclos `for`.

package main

import "fmt"

func main() {

    // El tipo más básico, con una condición sencilla.
    i := 1
    for i <= 3 {
        fmt.Println(i)
        i = i + 1
    }

    // El clásico ciclo `for` con estructura inicializar/condición/después.
    for j := 7; j <= 9; j++ {
        fmt.Println(j)
    }

    // `for` sin ninguna condición iterará repetidamente
    // hasta que se use `break` para salir del ciclo o `return`
    // para regresar un valor de la función que lo contiene.
    for {
        fmt.Println("loop")
        break
    }
}
