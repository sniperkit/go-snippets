#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static unsigned long next_permu(char *c, unsigned long strl) {
	char t;
	unsigned long i, k = strl, l = 0;

	for (i = strl - 1; i > 0; i--) {
		if (c[i - 1] < c[i]) {
			k = i - 1;
			break;
		}
	}
	if (k == strl)
		return 0;
	for (i = strl - 1; i > 0; i--) {
		if (c[i] > c[k]) {
			l = i;
			break;
		}
	}

	t = c[k];
	c[k] = c[l];
	c[l] = t;

	l = strl - k - 1;
	for (i = 0; (l - 1) - i > i; i++) {
		t = c[strl - 1 - i];
		c[strl - 1 - i] = c[k + 1 + i];
		c[k + 1 + i] = t;
	}
	return l;
}

static int ccmp(const void *b1, const void *b2) {
	char *c1 = (char *)b1, *c2 = (char *)b2;
	if (*c1 < *c2)
		return -1;
	return *c1 > *c2;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned long strl;
	char line[32];

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fgets(line, 31, fp) != 0) {
		strl = strlen(line);
		strl -= line[strl - 1] == '\n';
		char *c, *w = malloc(strl);
		c = strncpy(w, line, strl);
		qsort(c, strl, 1, ccmp);
		printf("%s", c);
		while (next_permu(c, strl))
			printf(",%s", c);
		printf("\n");
		free(w);
	}
	return 0;
}
