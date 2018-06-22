// Convertir cadenas numéricas a tipos númericos es una tarea básica y común en
// muchos programas, aquí vemos como hacerlo en Go.

package main

// El paquete `strconv` de la librería estándard contiene funciones auxiliares
// para convertir cadenas numéricas a [tipos numéricos de
// Go](http://golang.org/ref/spec#Numeric_types).
import "strconv"
import "fmt"

func main() {

    // El argumento `64` en éste caso de `strconv.ParseFloat` indica con
    // cuantos bits de precisión queremos convertir la cadena numérica a número
    // de punto flotante.
    f, _ := strconv.ParseFloat("1.234", 64)
    fmt.Println(f)

    // El argumento `0` en éste caso de `strconv.ParseInt` indica que la base
    // debe ser inferida directamente del formato de la cadena. El argumento
    // `64` hace que el resultado devuelto sea un número de punto flotante de
    // 64 bits.
    i, _ := strconv.ParseInt("123", 0, 64)
    fmt.Println(i)

    // La función `strconv.ParseInt` también puede reconocer números en formato
    // hexadecimal y convertirlos a entero.
    d, _ := strconv.ParseInt("0x1c8", 0, 64)
    fmt.Println(d)

    // Existe también una función `strconv.ParseUint`, útil para convertir una
    // cadena numérica a un número entero positivo (sin signo).
    u, _ := strconv.ParseUint("789", 0, 64)
    fmt.Println(u)

    // La función `strconv.Atoi` es una utilidad para conversión de cadena
    // numérica a entero.
    k, _ := strconv.Atoi("135")
    fmt.Println(k)

    // Las funciones de conversión regresan un error cuando la entrada no tiene
    // el formato numérico esperado.
    _, e := strconv.Atoi("wat")
    fmt.Println(e)
}
