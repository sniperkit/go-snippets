#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned sbs = 32, a, n = 0;
	char c;
	char *sb = malloc(sbs);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		if (c == '|') {
			while ((c = getc(fp)) == ' ') {
				fscanf(fp, "%u", &a);
				putchar(sb[a - 1]);
			}
			putchar('\n');
			n = 0;
		} else if (c != '\n') {
			if (n == sbs) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
			}
			sb[n++] = c;
		}
	}
	free(sb);
	return 0;
}
