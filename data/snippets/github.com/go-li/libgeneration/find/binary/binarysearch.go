// Package binary implement binary search
package binary

// Search searches function implemented and compared in calc, on interval
// [lo hi), until calc returns zero or avg returns nil. The inc adds
// a machine epsilon to a key.
// Search is parametrized by the key type, not by the type we are searching.
// Avg is the average function, if there is no more average it returns nil.
// Calc implements the custom calculation, and comparing the function values.
func Search(avg func(*, *) *, inc func(*), lo *, hi *, calc func(*) int) * {
	var mid *
	for {
		mid = avg(lo, hi)

		if mid == nil {
			break
		}

		var cmp = calc(mid)
		if cmp == 0 {
			return mid
		} else if cmp > 0 {
			if inc != nil {
				inc(mid)
			}
			lo = mid
		} else {
			hi = mid
		}
	}
	return hi
}

// Find finds a nearest key in integer domain [lo hi) that satisfies the find
// function. It optionally also searches a slice. If slice element for the
// current position does not exist, a nil is used.
// Search is parametrized by the type we are searching for (the slice element).
// Search returns the position of the equal or larger element in the slice.
// If the slice is decreasing and autodetection is on it finds equal or smaller.
func Find(autodetect bool, lo int, hi int, sli [], find func(*, int) int) int {
	var mul = 1
	if autodetect && (len(sli) > 0) && find(&sli[0], 0) < 0 {
		mul = -1
	}

	return *Search(func(low *int, high *int) *int {
		if *low == *high {
			return nil
		}
		var avg = (*low & *high) + ((*low ^ *high) >> 1)
		return &avg
	}, func(v *int) {
		(*v)++
	}, &lo, &hi, func(v *int) int {
		var x * = nil
		if (*v >= 0) && (*v < len(sli)) {
			x = &sli[*v]
		}
		return find(x, *v) * mul
	})
}
