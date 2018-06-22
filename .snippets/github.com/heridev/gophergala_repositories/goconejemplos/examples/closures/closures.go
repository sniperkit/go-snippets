// Go soporta _funciones anónimas_, con lo que se puede formar
// _closures_ o _cierres_.
// Las funciones anónimas son útiles cuando quieres definir una función
// a la cual no necesariamente quieres ponerle un nombre.

package main

import "fmt"

// Esta función `intSeq` regresa otra función, la cual definimos
// anónimanente en el cuerpo de `intSeq`. La función que regresamos
// _encierra_ a la variable `i` y esto forma un _closure_ o _cierre_.
func intSeq() func() int {
    i := 0
    return func() int {
        i += 1
        return i
    }
}

func main() {

    // Llamamos a `intSeq`, asignando el resultado (una función)
    // a la variable `nextInt`. El valor de esta función captura su
    // propio valor de `i`, el cual será actualizado cada vez que
    // llamamos a `nextInt`.
    nextInt := intSeq()

    // Aquí podemos ver el efecto del _cierre_ llamando varias veces
    // a `nextInt`.
    fmt.Println(nextInt())
    fmt.Println(nextInt())
    fmt.Println(nextInt())

    // Para confirmar que el estado de la variable es único a esa
    // función en particular creamos y probamos una nueva.
    newInts := intSeq()
    fmt.Println(newInts())
}
