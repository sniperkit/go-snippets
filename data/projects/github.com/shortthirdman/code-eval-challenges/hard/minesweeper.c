#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned m, n;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%u,%u;", &m, &n) != EOF) {
		unsigned d, i, j;
		int *prev = calloc(n, sizeof(int));
		int *curr = calloc(n, sizeof(int));

		for (j = 0; j < n; j++)
			if (fgetc(fp) == '*')
				prev[j] = -1;
			else
				prev[j] = (j > 0 && prev[j - 1] == -1);
		for (i = 1; i < m; i++) {
			if (fgetc(fp) == '*')
				curr[0] = -1;
			else
				curr[0] = (prev[0] == -1) + (n > 1 && prev[1] == -1);
			for (j = 1; j < n; j++) {
				if (fgetc(fp) == '*')
					curr[j] = -1;
				else
					curr[j] = (curr[j - 1] == -1) + (prev[j - 1] == -1) + (prev[j] == -1) + (n > j + 1 && prev[j + 1] == -1);
				if (prev[j - 1] == -1) {
					printf("*");
				} else {
					d = (unsigned)prev[j - 1] + (j > 1 && curr[j - 2] == -1) + (curr[j - 1] == -1) + (curr[j] == -1) + (prev[j] == -1);
					printf("%u", d);
				}
			}
			if (prev[j - 1] == -1) {
				printf("*");
			} else {
				d = (unsigned)prev[j - 1] + (j > 1 && curr[j - 2] == -1) + (curr[j - 1] == -1);
				printf("%u", d);
			}
			for (j = 0; j < n; j++)
				prev[j] = curr[j];
		}
		for (j = 0; j < n; j++) {
			if (prev[j] == -1) {
				printf("*");
			} else {
				d = (unsigned)prev[j] + (j < n - 1 && prev[j + 1] == -1);
				printf("%u", d);
			}
		}
		printf("\n");
		free(curr);
		free(prev);
	}
	return 0;
}
