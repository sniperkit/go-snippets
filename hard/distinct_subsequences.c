#include <stdio.h>
#include <stdlib.h>

static unsigned subs(char *p, char *q) {
	if (*q == '\0')
		return 1;
	if (*p == ',')
		return 0;
	if (*p == *q)
		return subs(p + 1, q + 1) + subs(p + 1, q);
	return subs(p + 1, q);
}

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned i = 0, sbs = 32;
	char c;
	char *sb = malloc(sbs), *q = NULL;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || i > 0) {
		if (c == '\n' || c == EOF) {
			sb[i] = '\0';
			printf("%u\n", subs(sb, q));
			i = 0;
		} else {
			if (i == sbs - 1) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
			}
			sb[i++] = c;
			if (c == ',')
				q = sb + i;
		}
	}
	free(sb);
	return 0;
}
