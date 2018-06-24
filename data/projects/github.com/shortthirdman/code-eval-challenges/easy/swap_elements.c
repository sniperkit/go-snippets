#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned ix = 0, ibs = 16, a, b, i, temp;
	int c;
	unsigned *ib = malloc(ibs * sizeof(unsigned));

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%u ", &a) != EOF) {
		ib[ix++] = a;
		while ((c = fscanf(fp, "%u ", &a)) > 0) {
			if (ix == ibs) {
				ibs += ibs / 2;
				ib = realloc(ib, ibs * sizeof(unsigned));
			}
			ib[ix++] = a;
		}
		fseek(fp, 1, SEEK_CUR);
		do {
			fscanf(fp, " %d-%d", &a, &b);
			temp = ib[a];
			ib[a] = ib[b];
			ib[b] = temp;
		} while ((c = getc(fp)) == ',');
		printf("%d", ib[0]);
		for (i = 1; i < ix; i++) {
			printf(" %d", ib[i]);
		}
		putchar('\n');
		ix = 0;
	}
	free(ib);
	return 0;
}
