#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned sbs = 64, w[11] = { 0 }, i = 0, j, k, m = 0, n;
	char c, *sb = malloc(sbs);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		if (c == '|') {
			fscanf(fp, " %u\n", &n);
			while (m > 1) {
				k = n % m;
				if (k > 0)
					for (j = k - 1; j < m - 1; j++)
						w[j] = w[j + 1];
				m--;
			}
			puts(sb + w[0]);
			i = 0;
			m = 0;
			w[0] = 0;
		} else {
			if (i == sbs) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
			}
			if (c == ' ') {
				sb[i++] = '\0';
				w[++m] = i;
			} else {
				sb[i++] = c;
			}
		}
	}
	free(sb);
	return 0;
}
