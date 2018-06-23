// Los _Canales_ son las tuberías que conectan
// goroutines concurrentes. Puedes enviar valores por
// un canal de una goroutine y recibir esos valores en
// otra goroutine.

package main

import "fmt"

func main() {

    // Puedes crear un canal nuevo con `make(chan val-type)`.
    // Los canales son del tipo de datos de los valores que
    // transportan.
    mensajes := make(chan string)

    // _Envía_ un valor por un canal usando la sintaxis
    // `canal <-`. Aquí estamos enviando `"ping"` al canal
    // de `mensajes` que creamos arriba, de otra nueva
    // goroutine.
    go func() { mensajes <- "ping" }()

    // La sintaxis `<-canal` _recibe_ un valor del canal especificado.
    // Aquí recibimos el mensaje `"ping"` que envíamos arriba
    // y lo mostramos en pantalla.
    msg := <-mensajes
    fmt.Println(msg)
}
