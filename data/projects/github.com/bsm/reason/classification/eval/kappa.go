package eval

import (
	"sync"

	"github.com/bsm/reason/core"
)

// Kappa represents the Cohen's kappa statistic.
type Kappa struct {
	m     [][]float64 // TODO: reimplement using gonum
	ncols int
	mu    sync.RWMutex
}

// NewKappa inits a new Kappa
func NewKappa() *Kappa {
	return new(Kappa)
}

// Record records an example of predicted vs actual with a given weight
func (k *Kappa) Record(predicted, actual core.Category, weight float64) {
	if !core.IsCat(predicted) || !core.IsCat(actual) {
		return
	}

	k.mu.Lock()
	defer k.mu.Unlock()

	if predicted > actual {
		k.grow(int(predicted + 1))
	} else {
		k.grow(int(actual + 1))
	}
	k.m[int(actual)][int(predicted)] += weight
}

// Score returns the kappa score
func (k *Kappa) Score() float64 {
	k.mu.RLock()
	defer k.mu.RUnlock()

	pws := make([]float64, len(k.m))
	tws := make([]float64, k.ncols)
	sum := 0.0

	for i, v := range k.m {
		for j, w := range v {
			pws[i] += w
			tws[j] += w
			sum += w
		}
	}
	if sum == 0.0 {
		return 0.0
	}

	var obs, exp float64
	for i, w := range pws {
		obs += k.m[i][i]
		exp += w * tws[i] / sum
	}
	if exp == sum {
		return 1.0
	}

	return (obs - exp) / (sum - exp)
}

func (k *Kappa) grow(n int) {
	if d := n - len(k.m); d > 0 {
		k.m = append(k.m, make([][]float64, d)...)
	}
	if k.ncols < n {
		for i, r := range k.m {
			if d := n - len(r); d > 0 {
				k.m[i] = append(r, make([]float64, d)...)
			}
		}
		k.ncols = n
	}
}
