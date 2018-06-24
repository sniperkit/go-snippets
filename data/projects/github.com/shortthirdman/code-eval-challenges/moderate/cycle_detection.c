#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

static bool uniq(unsigned *ib, unsigned s, unsigned i) {
	unsigned j;
	for (j = s + 1; j < i; j++)
		if (ib[s] == ib[j])
			return false;
	return true;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned a, s = 0, i = 0, ibs = 32;
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
			while (i - s > 1 && uniq(ib, s, i))
				s++;
			if (s < i - 1) {
				printf("%u", ib[s]);
				for (i = s + 1; ib[i] != ib[s]; i++)
					printf(" %u", ib[i]);
			}
			printf("\n");
			s = 0;
			i = 0;
		}
	}
	free(ib);
	return 0;
}
