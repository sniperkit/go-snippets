// Las [Variables de entorno](http://es.wikipedia.org/wiki/Variable_de_entorno)
// son un mecanismo universal para [transmitir datos de configuración a
// nuestros programas](http://www.12factor.net/config).
// A continuación veremos como definir, obtener y listar variables de entorno.

package main

import "os"
import "strings"
import "fmt"

func main() {

    // Para definir un par variable/valor usamos la función `os.Setenv`. Para
    // obtener el valor de una variable de entorno usamos `os.Getenv`, ésta
    // última función regresará una cadena vacía si la variable no está definida
    // en el entorno.
    os.Setenv("FOO", "1")
    fmt.Println("FOO:", os.Getenv("FOO"))
    fmt.Println("BAR:", os.Getenv("BAR"))

    // Podemos usar `os.Environ` para listar todos los pares variable/valor
    // presentes en el entorno. Ésta función regresa un slice de cadenas en la
    // forma `NOMBRE=valor`. Puedes usar `strings.Split` o `strings.SplitN` en
    // estos valores para separar el nombre de la variable y su valor. El
    // siguiente ejemplo  imprime el nombre de todas las variables que están
    // definidas en el entorno.
    fmt.Println()
    for _, e := range os.Environ() {
        pair := strings.Split(e, "=")
        fmt.Println(pair[0])
    }
}
