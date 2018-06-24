#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned i = 0, sbs = 32;
	char c;
	char *sb = malloc(sbs), *sbp = malloc(sbs);
	sbp[0] = '\0';

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || i > 0) {
		if (c == '\n' || c == EOF) {
			printf("%s\n", sbp);
			sbp[0] = '\0';
			i = 0;
		} else if (c == ' ') {
			char *temp = sbp;
			sb[i] = '\0';
			sbp = sb;
			sb = temp;
			i = 0;
		} else {
			if (i == sbs - 1) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
				sbp = realloc(sbp, sbs);
			}
			sb[i++] = c;
		}
	}
	free(sbp);
	free(sb);
	return 0;
}
