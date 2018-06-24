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

int main(void) {
	int primes[1000] = { 2, 3 };
	int i, c = 5, s = 5;
	for (i = 2; i < 1000; i++) {
		while (!prime(c, primes))
			c += 2;
		primes[i] = c;
		s += c;
		c += 2;
	}
	printf("%d\n", s);
	return 0;
}
