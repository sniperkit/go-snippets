package quick

var Seed uint32 = 1701105312
var Cutoff int = 20

func sort(a [], tmp *, cmp func(*, *) int) {
	for i := 1; i < len(a); i++ {
		var p = &a[i]
		*tmp = *p
		j := i - 1
		for ; j >= 0 && cmp(&a[j], tmp) > 0; j-- {
			var q = &a[j+1]
			var r = &a[j]
			*q = *r
		}
		p = &a[j+1]
		*p = *tmp
	}
}

func qsort(a [], tmp *, cmp func(*, *) int) {
	for {
		if len(a) < 2 {
			return
		}

		if len(a) < Cutoff {
			sort(a, tmp, cmp)
			return
		}

		left, right := 0, len(a)-1

		var random = uint64(uint32(uint64(Seed) * uint64(len(a))))
		pivotIndex := int((random * uint64(len(a))) >> 32)

		var p *
		var q *

		p = &a[pivotIndex]
		q = &a[right]

		*tmp = *p
		*p = *q
		*q = *tmp

		for i := range a {
			if cmp(&a[i], &a[right]) < 0 {

				p = &a[i]
				q = &a[left]

				*tmp = *p
				*p = *q
				*q = *tmp

				left++
			}
		}

		p = &a[right]
		q = &a[left]

		*tmp = *p
		*p = *q
		*q = *tmp

		qsort(a[0:left], tmp, cmp)
		a = a[left+1:]
	}
}

func Sort(a [], cmp func(*, *) int) {
	var temp *
	if len(a) >= cap(a) {
		var buf []
		buf = make([], 1)
		temp = &buf[0]
	} else {
		temp = &(a[0:cap(a)])[len(a)]
	}
	qsort(a, temp, cmp)
}
