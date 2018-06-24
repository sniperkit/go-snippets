package sparsedense

// BetterOffDense returns true if something should be dense.
func BetterOffDense(count, capacity int) bool {
	return count > 20 && count*12 > capacity
}
