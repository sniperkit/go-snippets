// En Go, un _arreglo_ es una secuencia numerada de elementos
// con un tamaño fijo.

package main

import "fmt"

func main() {

    // Aquí estamos creando un arreglo `a` que contendrá exactamente
    // 5 `int`s. El tipo de los elementos y el tamaño son parte del
    // tipo del arreglo. Por defecto un arreglo es de valor cero,
    // lo que para los valores de tipo `int` significa `0`s.
    var a [5]int
    fmt.Println("emp:", a)

    // Podemos establecer el valor de un elemento utilizando
    // la sintaxis `array[index] = value`, y obtener el valor
    // de un elemento utilizando `array[index]`.
    a[4] = 100
    fmt.Println("set:", a)
    fmt.Println("get:", a[4])

    // La sentencia `len` regresa el tamaño de un arreglo.
    fmt.Println("len:", len(a))

    // Usa esta sintaxis para declarar e inicializar un arreglo
    // en una línea.
    b := [5]int{1, 2, 3, 4, 5}
    fmt.Println("dcl:", b)

    // Los arreglos son unidimensionales, pero puedes componer
    // tipos para construir estructuras de datos multidimensionales.
    var twoD [2][3]int
    for i := 0; i < 2; i++ {
        for j := 0; j < 3; j++ {
            twoD[i][j] = i + j
        }
    }
    fmt.Println("2d: ", twoD)
}
