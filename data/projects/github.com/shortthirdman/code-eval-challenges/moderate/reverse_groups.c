#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int a;
	unsigned n, i = 0, ibs = 16, j, k, l;
	int *ib = malloc(ibs * sizeof(int));

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d", &a) != EOF) {
		if (i == ibs) {
			ibs += ibs / 2;
			ib = realloc(ib, ibs * sizeof(int));
		}
		ib[i++] = a;
		if (getc(fp) == ';') {
			j = 0, k = 0, l = 0;
			fscanf(fp, "%u", &n);
			while (j + n <= i) {
				for (k = j + n; k > j; k--) {
					printf("%d", ib[k - 1]);
					if (++l < i)
						printf(",");
				}
				j += n;
			}
			for (k = j; k < i; k++) {
				printf("%d", ib[k]);
				if (++l < i)
					printf(",");
			}
			printf("\n");
			fseek(fp, 1, SEEK_CUR);
			i = 0;
		}
	}
	free(ib);
	return 0;
}
