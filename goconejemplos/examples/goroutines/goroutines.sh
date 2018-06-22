# Cuando corremos este programa, podemos ver la salida
# de la llamada que bloquea primero, luego la salida
# intercalada de las dos goroutines. Esto refleja
# que las goroutines están siendo ejecutadas de manera
# concurrente por Go.
$ go run goroutines.go
direct : 0
direct : 1
direct : 2
goroutine : 0
going
goroutine : 1
goroutine : 2
<enter>
done

# Ahora veremos un complemento de las goroutine en
# programas concurrentes de Go: Canales.
