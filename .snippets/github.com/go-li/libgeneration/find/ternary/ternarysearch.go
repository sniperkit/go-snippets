// Package ternary implements ternary search
package ternary

// Search searches function implemented and compared in calc, on interval
// [lo hi), until calc returns zero or split3 returns nil as first result.
// Search is parametrized by the key type, not by the type we are searching.
// Split3 splits the current lo hi interval to 1/3 and 2/3, or it returns nil
// and the value in the middle (if the interval is too short).
// Calc implements the custom calculations, and comparing two function values.
// The function searched must be unimodal.
func Search(split3 func(*, *) (*, *), lo *, hi *, calc func(*, *) int) * {
	for {
		var newlo, newhi = split3(lo, hi)
		if newlo == nil {
			return newhi
		}

		if calc(newlo, newhi) > 0 {
			lo = newlo
		} else {
			hi = newhi
		}
	}
}

// Find finds a nearest key in integer domain [lo hi) that satisfies the find
// function. It optionally also searches a slice. If slice element for the
// current position does not exist, a nil is used, but index is also used.
// Find is parametrized by the type we are searching for (the slice element).
// Find returns the position of the global minimum.
// The function searched must be unimodal.
// If the slice contains duplicate elements, solution must be inbetween.
func Find(dt bool, lo int, hi int, sli [], find func(*, int, *, int) int) int {
	var mul = 1
	if dt && (len(sli) >= 2) && find(&sli[0], 0, &sli[1], 1) < 0 {
		mul = -1
	}
	return *Search(func(low *int, high *int) (*int, *int) {
		if *low+2 >= *high {
			var avg = (*low & *high) + ((*low ^ *high) >> 1)
			return nil, &avg
		}
		var lo = *low + (*high-*low)/3
		var hi = *low + ((*high-*low)*2)/3
		return &lo, &hi
	}, &lo, &hi, func(v *int, w *int) int {
		var x * = nil
		var y * = nil
		if (*v >= 0) && (*v < len(sli)) {
			x = &sli[*v]
		}
		if (*w >= 0) && (*w < len(sli)) {
			y = &sli[*w]
		}
		return find(x, *v, y, *w) * mul
	})
}
