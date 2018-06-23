// Una llamada a `panic` típicamente significa que sucedió
// un error inesperado. Lo usamos principalmente para terminar
// la ejecución en caso de errores que no debieran aparecer
// durante la operación normal o que no estamos listos
// para manejar adecuadamente.

package main

import "os"

func main() {

    // Usaremos `panic` a lo largo de este sitio para comprobar
    // errores inesperados. Este el el único programa en el
    // sitio diseñado para llamar a la función `panic`
    panic("un problema")

    // Un uso común de `panic` es abortar la ejecución
    // si una función devuelve un valor de error que no
    // sabemos (o queremos) manejar. Este es un ejemplo
    // de `panic` si encontramos un error inesperado al
    // momento de crear un archivo nuevo
    _, err := os.Create("/tmp/file")
    if err != nil {
        panic(err)
    }
}
