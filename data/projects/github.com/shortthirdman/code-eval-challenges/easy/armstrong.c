#include <stdio.h>

static int powi(int y, int x) {
	int r = 1;
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
	int n;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d;", &n) != EOF) {
		int a, e = 0, s = 0;
		for (a = n; a; a /= 10)
			e++;
		for (a = n; a; a /= 10)
			s += powi(a % 10, e);
		puts(n == s ? "True" : "False");
	}
	return 0;
}
