# Cuando corremos el programa el mensaje `"ping"` se
# pasa de una goroutine a otra a través de nuestro
# canal.
$ go run channels.go
ping

# Por defecto la recepción y los envíos se bloquean hasta
# que ambos receptor y transmisor están listos. Esta
# propiedad nos permite esperar hasta el mensaje
# `"ping"` hasta el final del programa sin tener
# que usar ningún otro tipo de sincronización.
