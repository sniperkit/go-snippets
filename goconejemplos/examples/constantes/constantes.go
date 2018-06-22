// Go soporta _constantes_ de tipo carácter, cadena,
// booleano y valores numéricos.

package main

import "fmt"
import "math"

// `const` se usa para declarar valores constantes.
const s string = "constant"

func main() {
    fmt.Println(s)

    // Una declaración `const` puede usarse donde mismo que
    // una declaración `var`.
    const n = 500000000

    // Las expresiones constantes se ejecutan con precisión
    // arbitraria.
    const d = 3e20 / n
    fmt.Println(d)

    // Una constante numérica no tiene tipo hasta que
    // se le asigna uno, por ejemplo con una conversión
    // explícita.
    fmt.Println(int64(d))

    // A un número se le puede dar el tipo utilizandolo
    // en un contexto que requiera uno, como la asignación
    // de una variable o la llamada a una función. Por ejemplo,
    // aquí `math.Sin` espera una variable de tipo `float64`.
    fmt.Println(math.Sin(n))
}
