// El lenguaje Go proporciona el paquete `math/rand`
// para la generaci&oacute;n de n&uacute;meros pseudo-aleatorios
// [n&uacute;mero pseudoaleatorio]
// (http://es.wikipedia.org/wiki/N%C3%BAmero_pseudoaleatorio).

package main

import "fmt"
import "math/rand"

func main() {
    // Por ejemplo, `rand.Intn` regresa un
    // n&uacute;mero aleatorio n de tipo `int`
    // dentro del rango: `0 <= n < 100`.
    fmt.Print(rand.Intn(100), ",")
    fmt.Print(rand.Intn(100))
    fmt.Println()

    // Por otro lado `rand.Float64` regresa
    // un n&uacute;mero `f` de tipo `float64`
    // dentro del rango:
    // `0.0 <= f < 1.0`.
    fmt.Println(rand.Float64())

    // El siguiente ejemplo se puede
    // usar para generar n&uacute;meros
    // aleatorios de punto flotante
    // en otros rangos, por ejemplo, el siguiente
    // c&oacute;digo genera valores
    // dentro del rango:
    // `5.0 <= f' < 10.0`.
    fmt.Print((rand.Float64()*5)+5, ",")
    fmt.Print((rand.Float64() * 5) + 5)
    fmt.Println()

    // Para generar un generador
    // pseudo-aleatorio determinista
    // es necesario iniciar con una semilla conocida.
    s1 := rand.NewSource(42)
    r1 := rand.New(s1)

    // Despu&eacute;s, s&oacute;lo hay que invocar al
    // resultado de `rand.Source`
    // igual que si invoc&aacute;ramos una
    // funci&oacute;n del paquete `rand`
    fmt.Print(r1.Intn(100), ",")
    fmt.Print(r1.Intn(100))
    fmt.Println()

    // Si proporcionas el mismo
    // valor semilla como entrada
    // se genera la misma secuencia
    // de n&uacute;meros aleatorios.
    s2 := rand.NewSource(42)
    r2 := rand.New(s2)
    fmt.Print(r2.Intn(100), ",")
    fmt.Print(r2.Intn(100))
    fmt.Println()
}
