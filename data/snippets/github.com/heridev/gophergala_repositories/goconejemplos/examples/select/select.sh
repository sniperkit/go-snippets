# Recibimos los valores `"uno"` y luego `"dos"`
# tal y como lo esperamos.
$ time go run select.go
recibido uno
recibido dos

# Nota que el tiempo de ejecución total es de
# tan solo ~2 segundos porque ambos `Sleeps`
# se ejecutan de manera concurrente
real  0m2.245s
