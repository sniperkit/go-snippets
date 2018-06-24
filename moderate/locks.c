#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned n, m, ls = 100, i, j, c;
	unsigned char *l = malloc(ls);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%u %u\n", &n, &m) != EOF) {
		c = n;
		if (m == 1) {
			c = n - 1;
		} else if (m) {
			if (n > ls) {
				ls = n;
				l = realloc(l, ls);
			}
			for (i = 0; i < n; i++)
				l[i] = 0;
			for (i = 0; i <= m % 2; i++) {
				for (j = 1; j < n; j += 2)
					l[j] = 1;
				for (j = 2; j < n; j += 3)
					if (l[j] != 0)
						l[j] = 0;
					else
						l[j] = 1;
			}
			if (l[n - 1] != 0)
				l[n - 1] = 0;
			else
				l[n - 1] = 1;
			for (i = 0; i < n; i++)
				c -= l[i];
		}
		printf("%u\n", c);
	}
	free(l);
	return 0;
}
