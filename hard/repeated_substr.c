#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

static bool allspace(char *p, unsigned n) {
	unsigned i;
	for (i = 0; i < n; i++)
		if (p[i] != ' ')
			return false;
	return true;
}

static bool contains(char *p, char *q, unsigned n, unsigned m) {
	unsigned i, j;
	for (i = 0; i < m; i++) {
		for (j = 0; j < n; j++)
			if (*(p + i + j) != *(q + j))
				break;
		if (j == n)
			return true;
	}
	return false;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	bool f;
	unsigned sbs = 32, s = 0, n, rs, i, m = 0;
	char c;
	char *sb = malloc(sbs);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || s > 0) {
		if (c == '\n' || c == EOF) {
			sb[s] = '\0';
			rs = 0;
			if (s > 0) {
				for (n = 1; n <= (s - rs) / 2; n++) {
					f = false;
					for (i = rs; i < s - n; i++) {
						if (!allspace(sb + i, n) && contains(sb + i + n, sb + i, n, s - n - i)) {
							m = i + n;
							rs = i;
							f = true;
							break;
						}
					}
					if (!f)
						break;
				}
			}
			if (m == 0) {
				puts("NONE");
			} else {
				sb[m] = '\0';
				puts(sb + rs);
				m = 0;
			}
			s = 0;
		} else {
			if (s == sbs - 1) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
			}
			sb[s++] = c;
		}
	}
	free(sb);
	return 0;
}
