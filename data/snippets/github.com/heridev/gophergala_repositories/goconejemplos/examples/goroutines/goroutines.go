// Una _goroutine_ es un hilo o thread de ejecución ligero.

package main

import "fmt"

func f(from string) {
    for i := 0; i < 3; i++ {
        fmt.Println(from, ":", i)
    }
}

func main() {

    // Supongamos que tenemos una llamada a la función
    // `f(s)`. Así es como la llamaríamos de la manera
    // tradicional, o de manera síncrona.
    f("direct")

    // Para llamar esta función en una goroutine, usamos
    // `go f(s)`. Esta nueva goroutine se ejecutara de
    // manera concurrente a la que la está llamando.
    go f("goroutine")

    // Puedes comenzar una goroutine con una llamada a
    // una función anónima.
    go func(msg string) {
        fmt.Println(msg)
    }("going")

    // Nuestras dos goroutines están corriendo de manera
    // asíncrona en dos goroutines separadas, así que
    // la ejecución avanza hasta aquí. Este código
    // `Scanln` requiere que presionemos una tecla antes
    // de que el programa termine.
    var input string
    fmt.Scanln(&input)
    fmt.Println("done")
}
