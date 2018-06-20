package fenwick

func bLSI(i int) int {
	return ((i) & -(i))
}

/* Sum: Return the sum of the first i elements, 0 through i-1. */
func Sum(a [], i int, add func(*,*,*)) * {
	var temporary = make([],1)
	var sum = &temporary[0]

	for ; i != 0; i -= bLSI(i) {
		add(sum, sum, & a[i-1] )
	}
	return sum
}

/* Add: Add delta to element i (and thus sum(a, j) for all j > i). */
func Add(a [], i int,  delta *, add func(*,*,*)) {
	for ; int(i) < len(a); i += bLSI(i+1) {
		add(&a[i],&a[i], delta)
	}
}

/* Optimize by avoiding terms that the two sums have in common
 * (more effective when range is small).
 */
func Range(a [], i, j int, add func(*,*,*), neg func(*)*) * {
	var temporary = make([],1)
	var sum = &temporary[0]

	for ; j > i; j -= bLSI(j) {
		add(sum, sum, & a[j-1] )
	}
	for ; i > j; i -= bLSI(i) {
		var n = neg(& a[i-1] )
		add(sum, sum, n)
		neg(& a[i-1] )
	}
	return sum
}

/* Get: Returns the value of the element at index i. */
func Get(a [], i int, add func(*,*,*), neg func(*)*) * {
	return Range(a, i, i+1, add, neg)
}



/* Set: Sets the value of the element at index i. */
func Set(a [], value *, i int, add func(*,*,*), neg func(*)*) {
	var temporary = make([],1)
	var sum = &temporary[0]

	add(sum, value , neg( Get(a, i, add, neg)))

	Add(a, i, sum, add);
}


func Init(a [], add func(*,*,*)) {
	var i, j int;
	for i = 0; int(i) < len(a); i++ {
		j = i + bLSI(i+1);
		if (int(j) < len(a)) {
			add(&a[j], &a[j], &a[i])
		}
	}
}
