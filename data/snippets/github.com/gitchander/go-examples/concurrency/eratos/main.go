package main

import "fmt"

func makeGenerator() <-chan int {
	c := make(chan int)
	go func() {
		for i := 2; ; i++ {
			c <- i
		}
	}()
	return c
}

func addFilter(in <-chan int, prime int) <-chan int {
	out := make(chan int)
	go func() {
		for {
			i := <-in
			if i%prime != 0 {
				out <- i
			}
		}
	}()
	return out
}

// Sieve of Eratosthenes
func main() {
	c := makeGenerator()
	for i := 0; i < 100; i++ {
		prime := <-c
		c = addFilter(c, prime)
		fmt.Println(prime)
	}
}
