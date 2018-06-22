// Go ofrece soporte integrado para [expresiones regulares](http://es.wikipedia.org/wiki/Expresi%C3%B3n_regular)
// Aquí hay algunos ejemplos de tareas comunes con expresiones regulares
// en Go.

package main

import "bytes"
import "fmt"
import "regexp"

func main() {

    // Esto prueba si un patrón coincide con una cadena.
    match, _ := regexp.MatchString("p([a-z]+)ch", "peach")
    fmt.Println(match)

    // En el ejemplo anterior usamos directamente un patrón cadena pero para
    // otras tareas necesitarás `compilar` una estructura optimizada
    // del tipo `RegExp`.
    r, _ := regexp.Compile("p([a-z]+)ch")

    // Hay varios métodos disponibles para estas estructuras.
    // Aquí se muestra una coincidencia con un patrón como
    // vimos anteriormente.
    fmt.Println(r.MatchString("peach"))

    // Aquí se encuentra la coincidencia con la expresión.
    fmt.Println(r.FindString("peach punch"))

    // Aquí también se encuentra la primer coincidencia pero
    // devuelve los índices de inicio y fin en lugar del
    // texto.
    fmt.Println(r.FindStringIndex("peach punch"))

    // Las variantes 'Submatch' incluyen información sobre las
    // coincidencias completas y parciales. Este ejemplo
    // devuelve información de las coincidencias de los patrones
    // `p([a-z]+)ch` and `([a-z]+)`.
    fmt.Println(r.FindStringSubmatch("peach punch"))

    // De igual manera este ejemplo devuelve información
    // sobre los índices de coincidencias y subcoincidencias.
    fmt.Println(r.FindStringSubmatchIndex("peach punch"))

    // Las variantes 'All' de las funciones aplican para todas
    // las coincidencias del texto de entrada, no solo para el
    // principio. Por ejemplo para encontrar todas las
    // coincidencias de una expresión regular.
    fmt.Println(r.FindAllString("peach punch pinch", -1))

    // Estas variantes 'All' están disponibles para las
    // demás funciones que vimos anteriormente.
    fmt.Println(r.FindAllStringSubmatchIndex(
        "peach punch pinch", -1))

    // Al pasar un entero no negativo en el segundo
    // parámetro limitará el numbero de coindencias.
    fmt.Println(r.FindAllString("peach punch pinch", 2))

    // Los ejemplos anteriores tienen argumentos de cadena
    // y usan nombres como 'MatchString'. También podemos
    // usar argumentos '[]byte' y eliminar 'String' del
    // nombre de la función.
    // Our examples above had string arguments and used
    // names like `MatchString`. We can also provide
    // `[]byte` arguments and drop `String` from the
    // function name.
    fmt.Println(r.Match([]byte("peach")))

    // Cuando creamos constantes con expresiones regulares
    // puedes utilizar la variante 'MustCompile' en lugar
    // de 'Compile'. La llamada 'Compile' no funcionará
    // para crear contantes porque devuelve dos valores.
    r = regexp.MustCompile("p([a-z]+)ch")
    fmt.Println(r)

    // El paquete `regexp` también puede ser usado para
    // reemplazar subcadenas de texto con otros valores.
    fmt.Println(r.ReplaceAllString("a peach", "<fruit>"))

    // La variante 'Func' permite transformar el texto
    // que coincide mediante una función.
    in := []byte("a peach")
    out := r.ReplaceAllFunc(in, bytes.ToUpper)
    fmt.Println(string(out))
}
