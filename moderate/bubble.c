#include <stdio.h>
#include <stdlib.h>

static unsigned min(unsigned a, unsigned b) {
	return a < b ? a : b;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	int c;
	unsigned a, i = 0, ibs = 32, j, k, x;
	unsigned *ib = malloc(ibs * sizeof(unsigned));

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%u ", &a) != EOF) {
		ib[i++] = a;
		while ((c = fscanf(fp, "%u ", &a)) > 0) {
			if (i == ibs) {
				ibs += ibs / 2;
				ib = realloc(ib, ibs * sizeof(unsigned));
			}
			ib[i++] = a;
		}
		fscanf(fp, "| %u", &a);
		a = min(a, i - 1);
		for (j = 0; j < a; j++) {
			for (k = 1; k < i; k++) {
				if (ib[k - 1] > ib[k]) {
					x = ib[k - 1];
					ib[k - 1] = ib[k];
					ib[k] = x;
				}
			}
		}
		for (j = 0; j < i; j++) {
			if (j)
				putchar(' ');
			printf("%u", ib[j]);
		}
		putchar('\n');
		i = 0;
	}
	free(ib);
	return 0;
}
