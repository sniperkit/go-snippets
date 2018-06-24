#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

static bool prime(unsigned a, unsigned *p) {
	while (*p * *p <= a) {
		if (a % *p == 0)
			return false;
		p++;
	}
	return true;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned i, n, last = 3, np = 2, ps = 8;
	unsigned *primes = malloc(ps * sizeof(unsigned));
	primes[0] = 2;
	primes[1] = 3;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%u", &n) != EOF) {
		while (last + 2 < n) {
			last += 2;
			if (prime(last, primes)) {
				if (ps == np) {
					ps += ps / 2;
					primes = realloc(primes, ps * sizeof(unsigned));
				}
				primes[np++] = last;
			}
		}
		if (n >= 2) {
			putchar('2');
			for (i = 1; i < np && primes[i] < n; i++)
				printf(",%u", primes[i]);
		}
		putchar('\n');
	}
	free(primes);
	return 0;
}
