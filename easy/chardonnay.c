#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

static int ccmp(const void *b1, const void *b2) {
	char *c1 = (char *)b1, *c2 = (char *)b2;
	if (*c1 < *c2)
		return -1;
	return *c1 > *c2;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned i, l, n;
	char wine[10][16], srtd[10][16], w[6];

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%s", &wine) != EOF) {
		bool f = false;
		unsigned m = 0;
		l = sprintf(srtd, "%s", wine);
		qsort(&srtd, l, 1, ccmp);
		for (n = 1; n < 10; n++) {
			fscanf(fp, "%s", wine[n]);
			if (wine[n][0] == '|')
				break;
			l = sprintf(srtd[n], "%s", wine[n]);
			qsort(&srtd[n], l, 1, ccmp);
		}
		if (n == 10)
			fseek(fp, 2, SEEK_CUR);
		fscanf(fp, "%s", &w);
		while (w[m] != '\0')
			m++;
		qsort(&w, m, 1, ccmp);

		for (i = 0; i < n; i++) {
			bool g = true;
			char *p = srtd[i], *q = w;
			while (*p != '\0' && *q != '\0') {
				if (*p < *q) {
					p++;
				} else if (*p == *q) {
					p++;
					q++;
				} else {
					g = false;
					break;
				}
			}
			if (g && *q == '\0') {
				if (f)
					putchar(' ');
				else
					f = true;
				printf("%s", wine[i]);
			}
		}
		if (!f)
			puts("False");
		else
			putchar('\n');
	}
	return 0;
}
