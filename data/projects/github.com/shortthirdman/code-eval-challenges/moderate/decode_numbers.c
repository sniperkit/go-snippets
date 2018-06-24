#include <stdio.h>
#include <stdlib.h>

static int decode(char *sb) {
	if (sb[0] == '\0' || sb[1] == '\0') {
		return 1;
	} else if (sb[0] != '1' && sb[0] != '2') {
		return decode(sb + 1);
	} else if (sb[1] == '0' || (sb[0] == '2' && sb[1] > '6')) {
		return decode(sb + 2);
	}
	return decode(sb + 2) + decode(sb + 1);
}

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned i = 0, sbs = 32;
	char c;
	char *sb = malloc(sbs);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || i > 0) {
		if (c == '\n' || c == EOF) {
			sb[i] = '\0';
			printf("%d\n", decode(sb));
			i = 0;
		} else {
			if (i == sbs - 1) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
			}
			sb[i++] = c;
		}
	}
	free(sb);
	return 0;
}
