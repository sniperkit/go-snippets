#include <stdio.h>
#include <stdlib.h>

static const char *alpha = " !\"#$%&'()*+,-./0123456789:<=>?@"
			   "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
			   "abcdefghijklmnopqrstuvwxyz";

static int indx(char c) {
	int i = 0;
	while (alpha[i] != c && i < 84)
		i++;
	return i;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned sbs = 8, s = 0, i = 0;
	char c, *sb = malloc(sbs);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		while (c != ';') {
			if (s == sbs) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
			}
			sb[s++] = c - '0';
			c = getc(fp);
		}
		while ((c = getc(fp)) != '\n' && c != EOF)
			printf("%c", alpha[(indx(c) - sb[i++ % s] + 84) % 84]);
		printf("\n");
		s = 0;
		i = 0;
	}
	free(sb);
	return 0;
}
