// La sentencia _switch_ expresa condicionales a través
// de varias derivaciones.

package main

import "fmt"
import "time"

func main() {

    // Ejemplo básico de `switch`.
    i := 2
    fmt.Print("write ", i, " as ")
    switch i {
    case 1:
        fmt.Println("one")
    case 2:
        fmt.Println("two")
    case 3:
        fmt.Println("three")
    }

    // Puedes usar comas para separar varias expresiones
    // en un mismo `case`. También usamos el caso `default`
    // opcional como parte de este ejemplo.
    switch time.Now().Weekday() {
    case time.Saturday, time.Sunday:
        fmt.Println("it's the weekend")
    default:
        fmt.Println("it's a weekday")
    }

    // Usar `switch` sin ninguna expresión es una manera alterna
    // de expresar lógica if/else. Aquí también mostramos como
    // las expresiones de la sentencia `case` no necesariamente deben
    // ser constantes.
    t := time.Now()
    switch {
    case t.Hour() < 12:
        fmt.Println("it's before noon")
    default:
        fmt.Println("it's after noon")
    }
}

// todo: switch con tipos.
