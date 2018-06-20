// Package binary implements binary heap
package binary

func percolateUp(array [], hole int, val *, cmp func(*, *) int) (next int) {
	for {
		next = ((hole + 1) >> 1) - 1
		if !((hole > 0) && (cmp(val, &(array)[next]) < 0)) {
			break
		}
		var l = &(array)[hole]
		var r = &(array)[next]
		*l = *r
		hole = next
	}
	return hole
}

func Insert(array *[], add [], cmp func(*, *) int) {
	if len(add) == 0 {
		return
	}
	if len(*array)+len(add) <= cap(*array) {
		*array = (*array)[0 : len(*array)+len(add)]
	} else if cap(*array) > 0 {
		*array = append(*array, add...)
	} else {
		*array = make([], len(add))
	}
	for i := range add {
		var q = &add[i]
		var p = &(*array)[percolateUp(*array, len(*array)-len(add)+i, q, cmp)]

		*p = *q
	}
}

func Heapify(array [], cmp func(*, *) int) {
	var temp *
	if len(array) >= cap(array) {
		var tmp []
		tmp = make([], 1)
		temp = &tmp[0]
	} else {
		temp = &(array[0 : len(array)+1])[len(array)]
	}
	for i := len(array)/2 - 1; i >= 0; i-- {
		down(array, i, temp, cmp)
	}
}

func down(array [], i0 int, tmp *, cmp func(*, *) int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= len(array) || j1 < 0 {
			break
		}
		j := j1
		if j2 := j1 + 1; j2 < len(array) && cmp(&array[j2], &array[j1]) < 0 {
			j = j2
		}
		if cmp(&array[j], &array[i]) > 0 {
			break
		}
		var pi *
		pi = &array[i]
		var pj *
		pj = &array[j]
		*tmp = *pi
		*pi = *pj
		*pj = *tmp
		i = j
	}
	return i > i0
}

func Fix(array [], i int, cmp func(*, *) int) {
	var temp *
	if len(array) >= cap(array) {
		var tmp []
		tmp = make([], 1)
		temp = &tmp[0]
	} else {
		temp = &(array[0 : len(array)+1])[len(array)]
	}
	if !down(array, i, temp, cmp) {
		percolateUp(array, i, &array[i], cmp)
	}
}

func Pop(array *[], cmp func(*, *) int) * {
	if len(*array) == 0 {
		return nil
	}
	var temp *
	if len(*array) >= cap(*array) {
		var tmp []
		tmp = make([], 1)
		temp = &tmp[0]
	} else {
		temp = &((*array)[0 : len(*array)+1])[len(*array)]
	}
	var p *
	p = &(*array)[0]
	var q *
	q = &(*array)[len(*array)-1]

	*temp = *p
	*p = *q
	*q = *temp

	*array = (*array)[0 : len(*array)-1]

	down(*array, 0, temp, cmp)

	return q
}

func Is(array [], cmp func(*, *) int) bool {
	for i := range array {
		if (len(array) > 2*i+1) && (cmp(&array[2*i+1], &array[i]) < 0) {
			return false
		}
		if (len(array) > 2*i+2) && (cmp(&array[2*i+2], &array[i]) < 0) {
			return false
		}
	}
	return true
}

func Peek(array [], _ func(*, *) int) * {
	if len(array) == 0 {
		return nil
	}
	return &array[0]
}

func Len(array []) int {
	return len(array)
}
