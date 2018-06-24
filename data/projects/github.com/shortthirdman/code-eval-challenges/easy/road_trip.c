#include <stdio.h>
#include <stdlib.h>

static int cmpint(const void *p1, const void *p2) {
	return *(int *)p1 - *(int *)p2;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	int a, b[600];
	unsigned i, n;
	char c;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	do {
		c = getc(fp);
	} while (c != ',' && c != EOF);
	while (fscanf(fp, "%d", &a) != EOF) {
		b[0] = a;
		n = 1;
		while (getc(fp) == ';' && getc(fp) == ' ') {
			do {
				c = getc(fp);
			} while (c != ',' && c != EOF);
			if (fscanf(fp, "%d", &b[n++]) == EOF)
				goto fail_eof;
		}
		qsort(&b, n, sizeof(int), cmpint);
		for (i = n - 1; i > 0; i--) {
			b[i] -= b[i - 1];
		}
		printf("%d", b[0]);
		for (i = 1; i < n; i++) {
			printf(",%d", b[i]);
		}
		printf("\n");
		do {
			c = getc(fp);
		} while (c != ',' && c != EOF);
	}
	return 0;
fail_eof:
	printf("unexpected end of input\n");
	return EXIT_FAILURE;
}
