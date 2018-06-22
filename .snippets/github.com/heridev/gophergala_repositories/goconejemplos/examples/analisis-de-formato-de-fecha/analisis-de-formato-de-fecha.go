// Go permite el dar formato y analisis de
// fecha mediante esquemas basados en patrones.

package main

import "fmt"
import "time"

func main() {
    p := fmt.Println

    // Aquí hay un ejemplo basico para dar formato
    // a la fecha de acuerdo a RFC3339, usando la
    // constante correspondiente en el esquema.

    t := time.Now()
    p(t.Format(time.RFC3339))

    // El analisis de la hora usa el mismo esquema de
    // valores como `Format`.

    t1, e := time.Parse(
        time.RFC3339,
        "2012-11-01T22:08:41+00:00")
    p(t1)

    // `Format` y `Parse` usan esquemas esquemas basados en
    // ejemplos. Normalmente usarás la constante `time` para
    // estos esquemas, aunque también puedes suministrar
    // esquemas personalizados. Los esquemas deben usar la hora
    // de referencia `Mon Jan 2 15:04:05 MST 2006` para mostrar
    // el patrón con el cual se puede dar formato/análisis una
    // determinada hora/cadena. El ejemplo de hora debe ser
    // exacamente como se muestra: el año 2006, 15 horas,
    // Lunes el día de la semana, etc.

    p(t.Format("3:04PM"))
    p(t.Format("Mon Jan _2 15:04:05 2006"))
    p(t.Format("2006-01-02T15:04:05.999999-07:00"))
    form := "3 04 PM"
    t2, e := time.Parse(form, "8 41 PM")
    p(t2)

    // Para representaciones exclusivamente numéricas
    // tambiés puedes hacer uso de una cadena el
    // formato estándar de cadena con los componentes
    // extraídos de el valor de la hora.

    fmt.Printf("%d-%02d-%02dT%02d:%02d:%02d-00:00\n",
        t.Year(), t.Month(), t.Day(),
        t.Hour(), t.Minute(), t.Second())

    // `Parse` regresará un error en entradas formadas
    // incorrectamente explicando el problema del análisis.

    ansic := "Mon Jan _2 15:04:05 2006"
    _, e = time.Parse(ansic, "8:41PM")
    p(e)
}
