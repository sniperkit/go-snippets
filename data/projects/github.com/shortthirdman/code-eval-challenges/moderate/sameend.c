#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned i = 0, j = 0, sbs = 32, k;
	char c;
	char *sb = malloc(sbs), *tb = malloc(sbs);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		if (c == ',') {
			bool f = true;
			while ((c = getc(fp)) != '\n' && c != EOF) {
				if (j > i)
					f = false;
				else
					tb[j++] = c;
			}
			if (f) {
				for (k = 1; k <= j; k++) {
					if (tb[j - k] != sb[i - k]) {
						f = false;
						break;
					}
				}
			}
			puts(f ? "1" : "0");
			i = 0;
			j = 0;
		} else if (c != '\n') {
			if (i == sbs) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
				tb = realloc(tb, sbs);
			}
			sb[i++] = c;
		}
	}
	free(sb);
	free(tb);
	return 0;
}
