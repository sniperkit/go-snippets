1. The commonly used asymptotic notation to calculate the runing time complexity of an algorithm:
	1. O Notation
	2. Ω Notation
	3. θ Notation

2. Big O Notation:
	The Notation O(n) is the formal way to express the upper bound of an algorithm's running time. It measures the worst case time complexity or the longest amount of time an algorithm can possibly take to complete.
	O(f(n)) = g(n) { there exists c > 0 and n0, so that g(n) <= c*f(n) for all n > n0. }
3. Omega Notation: 
	The Notation Ω(n) is the formal way to express the lower bound of an algorithm's runing time. It measures the best case complexity or the best amount of time an algorithm can possibly take to complete.
	Ω(f(n)) >= g(n) { there exist c > 0 and n0, so that g(n) <= c*f(n) for all n > n0. }
4. Theta Notation:
	The Notation θ(n) is the formal way to express both the lower bound and the upper bound of an algorithm's time. It is represented as follows:
	θ(f(n)) = g(n) { if and only if g(n) = O(f(n)) and g(n) = Ω(f(n)) for all n > n0. }
