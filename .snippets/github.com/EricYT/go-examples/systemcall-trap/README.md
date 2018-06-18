#systemcall in Go

1. `systemcall` trap by OS
`systemcall` like `file io` will block a thread of os in go,
when it was traped by OS. But if you need run other goroutine,
runtime will spawn new threads to execute goroutine.
Go don't wrap it into libaio, maybe it's not good choice for
a runtime of language.
