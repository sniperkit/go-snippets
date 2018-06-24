#include <stdio.h>
#include <stdlib.h>

static int cmpint(const void *p1, const void *p2) {
	return *(int *)p1 - *(int *)p2;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned ibs = 59, i = 0, k = 0;
	int *ib = malloc(ibs * sizeof(int));
	char c;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		if (c == '|') {
			qsort(ib, i, sizeof(int), cmpint);
			for (; i > 0; k = ib[k] % (1 << 16)) {
				putchar(ib[k] >> 16);
				i--;
			}
			putchar('\n');
			k = 0;
			fseek(fp, 1, SEEK_CUR);
			continue;
		}
		if (i == ibs) {
			ibs += ibs / 2;
			ib = realloc(ib, ibs * sizeof(int));
		}
		if (c == '$')
			k = i;
		ib[i] = (c << 16) + (int)i;
		i++;
	}
	free(ib);
	return 0;
}
