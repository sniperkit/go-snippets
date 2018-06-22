# Al correr el siguiente programa podemos observar que el
# valor asignado a la variable `FOO` dentro del programa
# es mostrado en pantalla, mientras que el valor de `BAR`
# queda vacío.
$ go run environment-variables.go
FOO: 1
BAR:

# La lista de variables de entorno disponibles depende del
# sistema operativo.
TERM_PROGRAM
PATH
SHELL
...

# Si definimos `BAR` en el entorno antes de llamar al
# programa, el último toma el valor del ambiente y lo
# imprime en pantalla.
$ BAR=2 go run environment-variables.go
FOO: 1
BAR: 2
...
