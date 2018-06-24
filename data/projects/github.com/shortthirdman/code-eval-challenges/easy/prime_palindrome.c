#include <stdbool.h>
#include <stdio.h>

static bool prime(int a, int *p) {
	while (*p * *p <= a) {
		if (a % *p == 0)
			return false;
		p++;
	}
	return true;
}

static bool pali(int x) {
	if ((x < 10) ||
	    (x < 100 && (x / 10 == x % 10)) ||
	    (x > 100 && (x / 100 == x % 10))) {
		return true;
	}
	return false;
}

int main(void) {
	int primes[256] = { 2, 3 };
	int i, m = 3, c = 2;
	for (i = 5; i < 1000; i += 2) {
		if (prime(i, primes)) {
			primes[c++] = i;
			if (pali(i))
				m = i;
		}
	}
	printf("%d\n", m);
	return 0;
}
