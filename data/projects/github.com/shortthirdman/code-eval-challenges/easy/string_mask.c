#include <stdbool.h>
#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	char c, s[20];
	int i = 0;
	bool p = true, w = true;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		if (w) {
			if (c == ' ') {
				i = 0;
				w = false;
			} else {
				s[i++] = c;
			}
		} else {
			if (c == '1') {
				putchar(s[i++] & 223);
			} else if (c == '0') {
				putchar(s[i++]);
			} else {
				putchar('\n');
				i = 0;
				w = true;
			}
			p = c == '\n';
		}
	}
	if (!p)
		putchar('\n');
	return 0;
}
