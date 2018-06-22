# Al correr nuestro programa se muestra una lista de
# cadenas ordenadas por longitud, como se esperaba.
$ go run sorting-by-functions.go
[kiwi peach banana]

# Siguiendo este mismo patrón de crear un tipo
# personalizado, implementar los 3 métodos de la
# `Interface` en ese tipo, y después llamar
# `sort.Sort` en una colección de ese tipo
# personalizado, podemos ordenar slices de Go
# mediante funciones arbitrarias.
