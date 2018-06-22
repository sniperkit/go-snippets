# Ejecutar este programa provocará la llamada a panic,
# imprimirá un mensaje de error, el registro de ejecución
# de la goroutine y terminará con un estatus diferente a
# cero.
$ go run panic.go
panic: un problema

goroutine 1 [running]:
main.main()
	/.../panic.go:12 +0x47
...
exit status 2

# Observe que a diferencia de algunos lenguajes que
# utilizan excepciones para el manejo de errores, en Go
# es propio del lenguaje devolver valores de error cada
# vez que sea posible
