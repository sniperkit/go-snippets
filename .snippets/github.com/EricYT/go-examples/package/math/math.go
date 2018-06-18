package math

func Average(xs []float64) float64 {
	total := 0.0
	for _, value := range xs {
		total += value
	}
	return total / float64(len(xs))
}
