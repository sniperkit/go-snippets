package splits

import "math"

// NormMerit normalises merit values
func NormMerit(m float64) float64 {
	if m < 0.0 || math.IsNaN(m) {
		return 0.0
	}
	return m
}

type GainRatioPenality struct {
	Weight float64
	value  float64
}

func (p *GainRatioPenality) Update(w float64) {
	rat := w / p.Weight
	p.value -= rat * math.Log2(rat)
}

func (p *GainRatioPenality) Value() float64 {
	if p.value <= 0.0 {
		return 1.0
	}
	return p.value
}
