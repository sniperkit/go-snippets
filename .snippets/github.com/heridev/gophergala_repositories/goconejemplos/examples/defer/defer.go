// _Defer_ se usa para asegurar que una función es
// llamada posteriormente durante la ejecución del
// programa, generalmente con propósitos de limpieza.
// `defer` se usa regularmente donde en otros lenguajes
// se utilizaría `ensure` y `finally`

package main

import "fmt"
import "os"

// Supongamos que queremos crear un archivo, escribir
// en él y luego cerrarlo al terminar. Así es como lo
// haríamos utilizando `defer`
func main() {

    // Inmediatamente después de obtener el objeto archivo
    // con `createFile`, diferimos el cierre del archivo con
    // `closeFile`. Esto se ejecutará al término de la función
    // contenedora (`main`), después de que `writeFile`
    // terminó de ejecutarse.
    f := createFile("/tmp/defer.txt")
    defer closeFile(f)
    writeFile(f)
}

func createFile(p string) *os.File {
    fmt.Println("crear")
    f, err := os.Create(p)
    if err != nil {
        panic(err)
    }
    return f
}

func writeFile(f *os.File) {
    fmt.Println("escribir")
    fmt.Fprintln(f, "data")

}

func closeFile(f *os.File) {
    fmt.Println("cerrar")
    f.Close()
}
