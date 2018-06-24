#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int l = 0, i = 0, j;
	char c, sb[10], cb[10];

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || i > 0) {
		if (c == ' ' || c == '\n' || c == EOF) {
			if (i > l)
				for (l = 0; l < i; l++)
					sb[l] = cb[l];
			if (c != ' ') {
				for (i = 0; i < l; i++) {
					if (i > 0)
						putchar(' ');
					for (j = 0; j < i; j++)
						putchar('*');
					putchar(sb[i]);
				}
				putchar('\n');
				l = 0;
			}
			i = 0;
		} else {
			i %= 10;
			cb[i++] = c;
		}
	}
	return 0;
}
