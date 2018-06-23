# Nota que aunque los slices son tipos de datos diferentes
# a los arreglos, se presentan de manera similar
# cuando se usa `fmt.Println`.
$ go run slices.go
emp: [  ]
set: [a b c]
get: c
len: 3
apd: [a b c d e f]
cpy: [a b c d e f]
sl1: [c d e]
sl2: [a b c d e]
sl3: [c d e f]
dcl: [g h i]
2d:  [[0] [1 2] [2 3 4]]

# Si quieres tener mas detalles del diseño y la
# implementación de los slices puedes leer
# este [excelente blog post](http://blog.golang.org/2011/01/go-slices-usage-and-internals.html)
# escrito por el equipo de Go.

# Ahora que hemos visto arreglos y slices vamos a ver
# la otra estructura básica de Go: Mapas.
