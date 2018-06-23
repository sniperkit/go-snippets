package devilc

type Range struct {
	Min, Max float64
	Count    int
}

func (r *Range) Step() float64 {
	return (r.Max - r.Min) / float64(r.Count-1)
}
