#include <stdbool.h>
#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	bool p = true;
	char c;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		if (c >= 'A' && c <= 'Z')
			c |= 32;
		putchar(c);
		p = c == '\n';
	}
	if (!p)
		putchar('\n');
	return 0;
}
