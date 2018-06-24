#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned sbs = 16, s = 0;
	char c;
	char *sb = malloc(sbs);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		if (c >= '0' && c <= '9') {
			sb[s++] = c;
			c = getc(fp);
			while ((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')) {
				if (s == sbs - 2) {
					sbs += sbs / 2;
					sb = realloc(sb, sbs);
				}
				sb[s++] = c;
				c = getc(fp);
			}
			if (c >= '0' && c <= '9') {
				sb[s] = sb[0];
				sb[0] = c;
			} else {
				sb[s] = c;
			}
			sb[s + 1] = '\0';
			printf("%s", sb);
		} else {
			putchar(c);
		}
		s = 0;
	}
	if (c != '\n')
		putchar('\n');
	free(sb);
	return 0;
}
