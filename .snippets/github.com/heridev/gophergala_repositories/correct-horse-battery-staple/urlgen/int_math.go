package urlgen

// Returns the numerator and denominator of a division operation
func divMod(numerator, denominator int64) (quotient, remainder int64) {
	return numerator / denominator, numerator % denominator
}

// Returns the integer result base^exp. Does not support negative values
func powInt(base, exp int64) int64 {
	if exp == 0 {
		return 1
	}

	storedBases := int64(1)

	for exp > 1 {
		if exp%2 == 1 {
			storedBases *= base
		}

		exp /= 2
		base *= base
	}

	return base * storedBases
}

// Returns the result of base^exp % mod. Does not support negative values
func powIntMod(base, exp, mod int64) int64 {
	if exp == 0 {
		return 1
	}

	storedBases := int64(1)

	for exp > 1 {
		if exp%2 == 1 {
			storedBases *= base
			storedBases %= mod
		}

		exp /= 2
		base *= base
		base %= mod
	}

	return (base * storedBases) % mod
}

// Returns the min of a and b
func intMin(a, b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func intSquareRoot(number int64) int64 {
	guess, oldGuess := number, number

	for {
		newGuess := (guess + (number / guess)) / 2

		if newGuess == guess {
			return guess
		}

		if oldGuess == newGuess {
			return intMin(guess, newGuess)
		}

		oldGuess = guess
		guess = newGuess
	}
}
