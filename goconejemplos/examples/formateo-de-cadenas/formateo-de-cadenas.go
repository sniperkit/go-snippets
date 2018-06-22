// Go ofrece un excelente soporte para el formato de 
// cadenas siguiendo la tradición de `printf`. Aquí hay
// algunos ejemplos de tareas comunes de formateo de
// cadenas.

package main

import "fmt"
import "os"

type point struct {
    x, y int
}

func main() {

    // Go ofrece varios "verbos" de impresión, diseñados
    // para dar formato a valores de Go simples. Por 
    // ejemplo, esto imprime una instancia de nuestra
    // estructura `point`.
    p := point{1, 2}
    fmt.Printf("%v\n", p)

    // Si el valor es una estructura, la varianete `%+v`
    // incluirá el nombre de los campos de la estructura.
    fmt.Printf("%+v\n", p)

    // La variante `%#v` imprime una representación de la
    // sintáxis en Go del valor, por ejemplo, el fragmento
    // de código que produciría ese valor.
    fmt.Printf("%#v\n", p)

    // Para imprimir el tipo de un valor, se usa `%T`.
    fmt.Printf("%T\n", p)

    // El formateo de boleanos es directo.
    fmt.Printf("%t\n", true)

    // Existen muchas opciones para formatear enteros.
    // Se usa `%d` para un formato base-10 estándar.
    fmt.Printf("%d\n", 123)

    // Esto imprime la representación binaria.
    fmt.Printf("%b\n", 14)

    // Esto imprime la letra que corresponda a ese entero
    fmt.Printf("%c\n", 33)

    // `%x` provee codificación hexadecimal.
    fmt.Printf("%x\n", 456)

    // Existen también varias opciones de formato para
    // números de punto flotante. Para formato decimal 
    // se usa `%f`.
    fmt.Printf("%f\n", 78.9)

    // `%e` y `%E` dan formato a los números de punto 
    // flotante usando (versiones ligeramente distintas
    //  de) la notación científica.
    fmt.Printf("%e\n", 123400000.0)
    fmt.Printf("%E\n", 123400000.0)

    // Para cadenas simples se usa `%s`.
    fmt.Printf("%s\n", "\"cadena\"")

    // Para incluir doble comilla como en el código Go
    // se usa `%q`.
    fmt.Printf("%q\n", "\"cadena\"")

    // Como con los enteros anteriormente, `%x` despliega
    // la cadena en base-16 usando dos letras en la 
    // salida por cada byte que haya en la entrada.
    fmt.Printf("%x\n", "hexadecimaleame esto")

    // Para imprimir la representación de un apuntador 
    // se usa `%p`.
    fmt.Printf("%p\n", &p)

    // Al dar formato a los números generalmente se desea
    // controlar el ancho y la precisión del resultado.
    // Para especificar el ancho de un entero, se usa un
    // número después del `%` en el verbo. Por omisión el
    // resultado estará justificado a la derecha usando
    // espacios.
    fmt.Printf("|%6d|%6d|\n", 12, 345)

    // También puedes especificar el ancho de los números
    // de punto flotante y generalmente también se quiere
    // restringir la precisión del punto decimal al mismo 
    // tiempo. Esto se logra usando la sintáxis:
    // ancho.precisión    
    fmt.Printf("|%6.2f|%6.2f|\n", 1.2, 3.45)

    // Para justificar a la izquierda se usa la bandera 
    // `-`.
    fmt.Printf("|%-6.2f|%-6.2f|\n", 1.2, 3.45)

    // También se puede querer controlar el ancho al dar
    // formato a cadenas, especialmente si se requiere 
    // que queden alineadas para salida tipo tabla. Para
    // justificación básica a la deerecha.
    fmt.Printf("|%6s|%6s|\n", "foo", "b")

    // Para justificar a la izquierda se usa la bandera 
    // `-` al igual que en los números.
    fmt.Printf("|%-6s|%-6s|\n", "foo", "b")

    // Hasta ahora hemos usado `Printf`, que imprime la
    // cadena formateada a `os.Stdout`. `Sprintf` le da
    // formato y regresa la cadena sin imprimirla en 
    // ningún lado.
    s := fmt.Sprintf("una %s", "cadena")
    fmt.Println(s)

    // Se puede formateo-imprimir a otros `io.Writers` 
    // además de `os.Stdout` usando `Fprintf`.
    fmt.Fprintf(os.Stderr, "un %s\n", "error")
}
