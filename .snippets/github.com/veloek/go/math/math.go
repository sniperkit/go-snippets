package math

// Collection of mathematical functions
//
// Made by Vegard Løkken <vegard@loekken.org>
// License MIT

// IsPrime checks wether or not a number is a prime
// It returns a boolean value of true or false
func IsPrime(n int) bool {
	if n <= 1 {
		return false
	} else if n <= 3 {
		return true
	} else if n%2 == 0 || n%3 == 0 {
		return false
	}

	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}

	return true
}
