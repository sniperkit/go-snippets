// Go soporta definir _métodos_ para structs.

package main

import "fmt"

type rect struct {
    width, height int
}

// Este método `area` tiene un _tipo receptor_ `*rect`.
func (r *rect) area() int {
    return r.width * r.height
}

// Los métodos pueden ser definidos para receptores de tipo
// apuntador o por valor. Aquí un ejemplo de un receptor
// por valor.
func (r rect) perim() int {
    return 2*r.width + 2*r.height
}

func main() {
    r := rect{width: 10, height: 5}

    // Aquí llamamos a los dos métodos definidos para nuestro struct.
    fmt.Println("area: ", r.area())
    fmt.Println("perim:", r.perim())

    // Go automáticamente maneja la conversión entre valores
    // y apuntadores en las llamadas a métodos. Puede que
    // quieras utilizar un tipo apuntador receptor para evitar copiar
    // en las llamadas a métodos o permitir que el método cambie
    // el struct receptor.
    rp := &r
    fmt.Println("area: ", rp.area())
    fmt.Println("perim:", rp.perim())
}
