// En ocasiones queremos ordenar una colección por algo
// diferente a su orden natural. Por ejemplo, supongamos que
// queremos ordenar cadenas por su longitud en vez de
// en orden alfabético. Aquí hay un ejemplo de un ordenamiento
// personalizado en Go.

package main

import "sort"
import "fmt"

// Para ordenar mediante una funcion personalizada en Go, necesitamos
// un tipo correspondiente. Aquí hemos creado un tipo llamado `ByLength`
// que solo es un alias para el tipo primitivo `[]string`
type ByLength []string

// Implementamos la interfaz `sort.Interface` - `Len`, `Less`, y
// `Swap` - en nuestro tipo, para que  podamos usar la función `Sort` genérica
// del paquete `sort`. `Len` y `Swap`
// usualmente serán similares entre tipos y `Less`
// sostendrá la lógica del ordenamiento personalizado. En nuestro caso
// queremos ordenar en base a la longitud ascendente de la cadena
// así que aquí usaremos `len(s[i])` y `len(s[j])`.
func (s ByLength) Len() int {
    return len(s)
}
func (s ByLength) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s ByLength) Less(i, j int) bool {
    return len(s[i]) < len(s[j])
}

// Con todo esto en su lugar, ahora podemos implementar nuestro
// ordenamiento personalizado al pasar como parámetro el slice original `fruits` en
// `ByLength`, y después usar la funcion `sort.Sort` en el slice tipeado.
func main() {
    fruits := []string{"peach", "banana", "kiwi"}
    sort.Sort(ByLength(fruits))
    fmt.Println(fruits)
}
