#include <stdio.h>

static int dsig(int b) {
	int d = 0, r;
	while (b) {
		r = b % 10;
		if (r)
			d += 1 << (3 * r);
		b /= 10;
	}
	return d;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	int a;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d", &a) != EOF) {
		int b = a + 9, d = dsig(a);
		while (d != dsig(b))
			b += 9;
		printf("%d\n", b);
	}
	return 0;
}
