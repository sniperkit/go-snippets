package main

type Interface interface{}

type Struct struct{}

func main() {
	var ps *Struct
	var pi *Interface
	pi = new(Interface)
	*pi = ps

	_, _ = pi, ps
}
