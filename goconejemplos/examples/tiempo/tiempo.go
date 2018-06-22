// Go nos ofrece un amplio soporte para tiempos y duraciones;
// aquí hay algunos ejemplos.

package main

import "fmt"
import "time"

func main() {
    p := fmt.Println

    // Comenzaremos obteniendo la hora actual.
    now := time.Now()
    p(now)

    // Puedes construir una estructura `time` al proveerle el
    // año, mes, día, etc. Las horas siempre están asociadas
    // con una `Location`, i.e. zona horaria.
    then := time.Date(
        2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
    p(then)

    // Puedes extraer los diferentes componentes de un valor de
    // tiempo tal y como lo esperabas.
    p(then.Year())
    p(then.Month())
    p(then.Day())
    p(then.Hour())
    p(then.Minute())
    p(then.Second())
    p(then.Nanosecond())
    p(then.Location())

    // Los días de la semana, `Weekday` Lunes-Domingo, también estan disponibles.
    p(then.Weekday())

    // Estos métodos comparan dos tiempos, verificando si
    // el primero ocurren antes, despues, o al mismo tiempo
    // que el segundo, respectivamente.
    p(then.Before(now))
    p(then.After(now))
    p(then.Equal(now))

    // El método `Sub` nos regresa una `Duration` representando
    // el intervalo entre dos tiempos.
    diff := now.Sub(then)
    p(diff)

    // Podemos computar la longitud de la duración
    // en diferentes unidades.
    p(diff.Hours())
    p(diff.Minutes())
    p(diff.Seconds())
    p(diff.Nanoseconds())

    // Puedes usar `Add` para incrementar un tiempo, mediante
    // una determinada duración, o con un `-` para
    // decrementarla.
    p(then.Add(diff))
    p(then.Add(-diff))
}
