#include <stdio.h>
#include <gmp.h>

static unsigned long power(unsigned long y, unsigned long x) {
	unsigned long r = 1;
	while (x) {
		if (x & 1)
			r *= y;
		y *= y;
		x >>= 1;
	}
	return r;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned long m, i;
	mpz_t r, t1, t2;
	mpz_init(r);
	mpz_init(t1);
	mpz_init(t2);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%lu", &m) != EOF) {
		mpz_set_ui(r, 0);
		for (i = 0; i < m / 2; i++) {
			mpz_bin_uiui(t1, m, i);
			mpz_bin_uiui(t2, 11 * m / 2 - 1 - 10 * i, m - 1);
			mpz_mul(t1, t1, t2);
			mpz_mul_si(t1, t1, power(-1, i));
			mpz_add(r, r, t1);
		}
		mpz_out_str(NULL, 10, r);
		putchar('\n');
	}
	mpz_clear(r);
	mpz_clear(t1);
	mpz_clear(t2);
	return 0;
}
