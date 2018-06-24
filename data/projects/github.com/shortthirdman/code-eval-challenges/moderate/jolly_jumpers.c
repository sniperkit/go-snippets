#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned n, m = 0, i, x;
	int xo, xc;
	bool *diff = NULL;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d", &n) != EOF) {
		bool jolly = true, done = false;
		if (diff == NULL) {
			diff = calloc(n, sizeof(bool));
			m = n;
		} else {
			if (m < n) {
				diff = realloc(diff, n * sizeof(bool));
				m = n;
			}
			for (i = 0; i < n; i++)
				diff[i] = false;
		}
		if (n == 1)
			done = true;
		fscanf(fp, "%d", &xo);
		for (i = 1; i < n; i++) {
			fscanf(fp, "%d", &xc);
			if (done)
				continue;
			x = (unsigned)abs(xc - xo);
			if (x >= n || x == 0 || diff[x - 1]) {
				jolly = false;
				done = true;
			} else {
				diff[x - 1] = true;
				xo = xc;
			}
		}
		puts(jolly ? "Jolly" : "Not jolly");
	}
	free(diff);
	return 0;
}
