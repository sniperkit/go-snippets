#include <stdbool.h>
#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	char c;
	bool p = true, u = false;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		if ((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			u = !u;
			if (u)
				putchar(c & 223);
			else
				putchar(c | 32);
		} else {
			if (c == '\n')
				u = false;
			putchar(c);
		}
		p = c == '\n';
	}
	if (!p)
		putchar('\n');
	return 0;
}
