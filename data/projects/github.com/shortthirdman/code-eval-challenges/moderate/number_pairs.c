#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned a, i = 0, ibs = 32, j, k, n;
	unsigned *ib = malloc(ibs * sizeof(unsigned));
	char c;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d", &a) != EOF) {
		if (i == ibs) {
			ibs += ibs / 2;
			ib = realloc(ib, ibs * sizeof(unsigned));
		}
		ib[i++] = a;
		if ((c = getc(fp)) == ';') {
			fscanf(fp, "%d", &a);
			n = 0;
			for (j = 0; j < i - 1 && 2 * ib[j] < a; j++)
				for (k = i - 1; k > j && ib[j] + ib[k] >= a; k--) {
					i--;
					if (ib[j] + ib[k] == a) {
						if (n++ > 0)
							putchar(';');
						printf("%d,%d", ib[j], ib[k]);
					}
				}
			if (n == 0)
				printf("NULL");
			putchar('\n');
			i = 0;
		}
	}
	free(ib);
	return 0;
}
