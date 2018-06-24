#include <inttypes.h>
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
	unsigned psize = 6, plast = 1, pmax = 3, i;
	uint64_t a;
	unsigned *primes = malloc(psize * sizeof(unsigned));
	primes[0] = 2;
	primes[1] = 3;
	for (i = primes[plast] + 2; plast < psize - 1; i += 2)
		if (prime(i, primes))
			primes[++plast] = i;
	pmax = i - 2;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%"SCNu64, &a) != EOF) {
		while (primes[plast] < 64 && a > pmax) {
			psize += psize / 2;
			primes = realloc(primes, psize * sizeof(unsigned));
			for (i = primes[plast] + 2; plast < psize - 1; i += 2)
				if (prime(i, primes))
					primes[++plast] = i;
			pmax = i - 2;
		}
		for (i = 0; primes[i] < 64 && ((uint64_t)1 << primes[i]) - 1 < a; i++) {
			if (i > 0)
				printf(", ");
			printf("%"PRIu64, ((uint64_t)1 << primes[i]) - 1);
		}
		printf("\n");
	}
	free(primes);
	return 0;
}
