package urlgen

var primes []int64

func init() {
	primes = append(primes, 2)
	primes = append(primes, 3)
}

func isPrime(number int64) bool {
	maxTest := intSquareRoot(number)
	testPrimes := getPrimesUpTo(maxTest)
	for _, prime := range testPrimes {
		if number%prime == 0 {
			return false
		}
	}

	return true
}

func getPrimeFactors(number int64) []int64 {
	maxTest := intSquareRoot(number)
	testPrimes := getPrimesUpTo(maxTest)
	factors := make([]int64, 0)
	for _, prime := range testPrimes {
		for number%prime == 0 {
			number /= prime
			factors = append(factors, prime)
		}
		if number == 1 {
			break
		}
	}

	if number != 1 {
		factors = append(factors, number)
	}

	return factors
}

func getUniquePrimeFactors(number int64) []int64 {
	maxTest := intSquareRoot(number)
	testPrimes := getPrimesUpTo(maxTest)
	factors := make([]int64, 0)
	for _, prime := range testPrimes {
		if number%prime == 0 {
			factors = append(factors, prime)
		}

		for number%prime == 0 {
			number /= prime
		}
		if number == 1 {
			break
		}
	}

	if number != 1 {
		factors = append(factors, number)
	}

	return factors
}

func getPrimesUpTo(number int64) []int64 {
	lastPrime := primes[len(primes)-1]
	if lastPrime >= number {
		// Binary search for the index
		bottom, top := 0, len(primes)-1
		guess := (top + bottom) / 2

		for {
			guessPrime := primes[guess]
			if guessPrime <= number {
				if primes[guess+1] > number {
					return primes[:guess+1]
				} else if primes[guess+1] == number {
					return primes[:guess]
				}
			}

			if guessPrime > number {
				top = guess
			} else {
				bottom = guess
			}

			guess = (top + bottom) / 2
		}
	}

	for lastPrime <= number {
		canidate := lastPrime + 2
		canidateSquareRoot := intSquareRoot(canidate)
		foundPrime := true

		for _, prime := range primes {
			if prime > canidateSquareRoot {
				break
			}

			if canidate%prime == 0 {
				foundPrime = false
				break
			}
		}

		if foundPrime {
			primes = append(primes, canidate)
		}

		lastPrime = canidate
	}

	return primes[:len(primes)]
}
