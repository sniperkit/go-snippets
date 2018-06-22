// La derivación con `if` y `else` en Go
// se hace de manera directa.

package main

import "fmt"

func main() {

    // Ejemplo básico.
    if 7%2 == 0 {
        fmt.Println("7 es par")
    } else {
        fmt.Println("7 es impar")
    }

    // Puedes utilizar un `if` sin un else.
    if 8%4 == 0 {
        fmt.Println("8 es divisible entre 4")
    }

    // Los condicionales pueden ser precedidos por
    // una declaración; cualquier
    // variable declarada en dicha declaración estará disponible
    // en todas las derivaciones.
    if num := 9; num < 0 {
        fmt.Println(num, "es negativo")
    } else if num < 10 {
        fmt.Println(num, "tiene 1 digito")
    } else {
        fmt.Println(num, "tiene multiples digitos")
    }
}

// Nota que no necesitas los paréntesis alrededor de las
// condiciones en Go, pero las llaves {} si son obligatorias.
