#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned sbs = 32, l = 0, i = 0;
	char c;
	char *sb = malloc(sbs), *cb = malloc(sbs);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || i > 0) {
		if (c == ' ' || c == '\n' || c == EOF) {
			if (i > l) {
				char *tb = sb;
				sb = cb;
				cb = tb;
				sb[i] = '\0';
				l = i;
			}
			i = 0;
			if (c != ' ') {
				puts(sb);
				l = 0;
			}
		} else {
			if (i == sbs - 1) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
				cb = realloc(cb, sbs);
			}
			cb[i++] = c;
		}
	}
	free(cb);
	free(sb);
	return 0;
}
