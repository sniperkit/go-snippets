// La librería estándar del paquete `strings` provee muchas
// funciones útiles relacionadas a las cadenas. Aquí hay algunos ejemplos
// donde puedes obtener noción del paquete.

package main

import s "strings"
import "fmt"

// Creamos un alias de `fmt.Println` a un nombre más corto pues
// lo usaremos bastante debajo.
var p = fmt.Println

func main() {

    // Aquí hay un ejemplo de las funciones disponibles en
    // `strings`. Nota que todas estas funciones son del
    // paquete, por lo que no son métodos del objeto string.
    // Esto significa que necesitaremos pasar la cadena en cuestión
    // como el primer parámetro de la función.
    p("Contains:  ", s.Contains("test", "es"))
    p("Count:     ", s.Count("test", "t"))
    p("HasPrefix: ", s.HasPrefix("test", "te"))
    p("HasSuffix: ", s.HasSuffix("test", "st"))
    p("Index:     ", s.Index("test", "e"))
    p("Join:      ", s.Join([]string{"a", "b"}, "-"))
    p("Repeat:    ", s.Repeat("a", 5))
    p("Replace:   ", s.Replace("foo", "o", "0", -1))
    p("Replace:   ", s.Replace("foo", "o", "0", 1))
    p("Split:     ", s.Split("a-b-c-d-e", "-"))
    p("ToLower:   ", s.ToLower("TEST"))
    p("ToUpper:   ", s.ToUpper("test"))
    p()

    // Puedes encontrar más funciones en la documentación del paquete [`strings`](http://golang.org/pkg/strings/)

    // No son parte de `strings` pero valen la pena mencionar aquí
    // los mecanismos para obtener la longitud de una cadena
    // y para obtener un carácter por índice.
    p("Len: ", len("hello"))
    p("Char:", "hello"[1])
}
