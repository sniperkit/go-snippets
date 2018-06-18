package main

import "fmt"

func main() {
	xs := []float64{12, 23.23, 234, 342, 23.2}

	average := average(xs)
	fmt.Println("average is ", average)

	fmt.Println(add(1, 2, 3))

	// even generator
	evenGenerator := makeEvenGenerator()
	fmt.Println(evenGenerator())
	fmt.Println(evenGenerator())
	fmt.Println(evenGenerator())

	// factorial
	fmt.Println(factorial(10))
}

/* average function */
func average(xs []float64) float64 {
	total := 0.0
	for _, value := range xs {
		total += value
	}
	return total / float64(len(xs))
}

/* sum */
func add(args ...int) (sum int) {
	for _, value := range args {
		sum += value
	}
	return
}

/* even generate */
func makeEvenGenerator() func() uint {
	i := uint(0)
	return func() (r uint) {
		r = i
		i += 2
		return
	}
}

/* factorial */
func factorial(x uint) uint {
	if x == 0 {
		return 1
	}
	return x * factorial(x-1)
}
