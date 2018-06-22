// [_Command-line arguments_](http://en.wikipedia.org/wiki/Command-line_interface#Arguments)
// Los argumendos de la l&iacute;nea de comandos son una
// forma com&uacute;n para parametrizar la ejecuci&oacute;n de programas
// Por ejemplo, `go run hello.go` toma `run` y `hello.go`
// como argumentos para el ejecutable `go`

package main

import "os"
import "fmt"

func main() {

    // La sentencia `os.Args` proporciona
    // acceso directo a los argumentos
    // pasados por la l&iacute;nea de comandos.
    // Observa que el primer elemento en
    // esta secuencia es la ruta hacia el programa.
    // Y `os.Args[1:]`  contiene los argumentos del programa
    argsWithProg := os.Args
    argsWithoutProg := os.Args[1:]

    // Puedes  obtener tambien de manera individual los argumentos
    // usando un acceso basado en indices.
    arg := os.Args[3]

    fmt.Println(argsWithProg)
    fmt.Println(argsWithoutProg)
    fmt.Println(arg)
}
