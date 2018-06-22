// En Go, las _variables_ son declaradas explícitamente y
// usadas por el compilador para e.g. asegurar que las
// llamadas a las funciones sean del tipo correcto.

package main

import "fmt"

func main() {

    // `var` se usa para declarar una o más variables.
    var a string = "initial"
    fmt.Println(a)

    // Puedes declarar múltiples variables en una línea.
    var b, c int = 1, 2
    fmt.Println(b, c)

    // Go infiere el tipo de las variables inicializadas.
    var d = true
    fmt.Println(d)

    // Las variables declaradas sin su inicialización correspondiente
    // son de _valor-cero_. Por ejemplo, el valor cero de una
    // varibale de tipo `int` es `0`.
    var e int
    fmt.Println(e)

    // La syntaxis `:=` es la abreviación para declarar e inicializar
    // una variable, e.g. de `var f string = "short"` en este caso.
    f := "short"
    fmt.Println(f)
}
