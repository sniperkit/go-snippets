#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

static int shrink(int a) {
	return 4 * a / 5;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned a, i = 0, ibs = 32, j, t;
	unsigned *ib = malloc(ibs * sizeof(unsigned));
	char c;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%u", &a) != EOF) {
		if (i == ibs) {
			ibs += ibs / 2;
			ib = realloc(ib, ibs * sizeof(unsigned));
		}
		ib[i++] = a;
		if ((c = getc(fp)) == '\n' || c == EOF) {
			unsigned n = 0, gap = i, last = 0;
			bool changed = false;
			while (gap > 1 || changed) {
				changed = false;
				n++;
				if (gap > 1)
					gap = shrink(gap);
				for (j = 0; j < i - gap; j++)
					if (ib[j] > ib[j + gap]) {
						t = ib[j];
						ib[j] = ib[j + gap];
						ib[j + gap] = t;
						changed = true;
						last = n;
					}
			}
			printf("%u\n", last);
			i = 0;
		}
	}
	free(ib);
	return 0;
}
