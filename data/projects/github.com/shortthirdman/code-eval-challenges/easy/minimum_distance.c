#include <stdio.h>
#include <stdlib.h>

static int cmpint(const void *p1, const void *p2) {
	return *(int *)p1 - *(int *)p2;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned i, n, f[99], r, s;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%u ", &n) != EOF) {
		r = 0;
		for (i = 0; i < n; i++)
			fscanf(fp, "%u ", &f[i]);
		qsort(&f, n, sizeof(unsigned), cmpint);
		s = f[n / 2];
		for (i = 0; i < n; i++)
			r += f[i] > s ? f[i] - s : s - f[i];
		printf("%d\n", r);
	}
	return 0;
}
