#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>

static int pow3(int x) {
	int y = 3, r = 1;
	while (x) {
		if (x & 1)
			r *= y;
		y *= y;
		x >>= 1;
	}
	return r;
}

static bool ugly(int j, int i, int *n) {
	int k = 0, cj = j;
	int64_t s = 0, c = n[0];
	bool p = true;

	while (k < i - 1) {
		int ops = cj % 3;
		cj /= 3;
		if (ops == 0) {
			c *= 10;
		} else {
			if (p) {
				s += c;
				if (ops == 1)
					p = false;
			} else {
				s -= c;
				if (ops == 2)
					p = true;
			}
			c = 0;
		}
		k++;
		c += n[k];
	}
	if (p)
		s += c;
	else
		s -= c;
	if (s % 2 == 0 || s % 3 == 0 || s % 5 == 0 || s % 7 == 0)
		return true;
	return false;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	char c;
	int j, n[13];

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		int i = 0, u = 0;
		while (c != '\n' && c != EOF) {
			n[i++] = c - '0';
			c = getc(fp);
		}
		for (j = pow3(i - 1); j > 0 ; j--)
			if (ugly(j - 1, i, n))
				u++;

		printf("%d\n", u);
	}
	return 0;
}
