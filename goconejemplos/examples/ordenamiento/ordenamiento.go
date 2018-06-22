// El paquete `sort` de Go implementa ordenamiento para los tipos primitivos
// y definidos por el usuario. Primero veremos el ordenamiento para los
// del tipo primitivo.

package main

import "fmt"
import "sort"

func main() {

    // Los métodos de ordenamiento son específicos para los tipos primitivos;
    // aquí hay un ejemplo para `strings`. Observa que el ordenamiento es
    // interno, por lo que se cambia el slice proporcionado y no se
    // regresa uno nuevo.
    strs := []string{"c", "a", "b"}
    sort.Strings(strs)
    fmt.Println("Strings:", strs)

    // Un ejemplo de ordenamiento de `int`s.
    ints := []int{7, 2, 4}
    sort.Ints(ints)
    fmt.Println("Ints:   ", ints)

    // También podemos usar `sort` para verificar si un slice
    // ya está actualmente ordenado.
    s := sort.IntsAreSorted(ints)
    fmt.Println("Sorted: ", s)
}
