package main

func main() {

	fff := func() {

	}

	fff()

	var ccc func()
	ccc = func() {
		ccc()
	}

	ccc()
}