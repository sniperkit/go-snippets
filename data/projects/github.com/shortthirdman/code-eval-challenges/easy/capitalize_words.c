#include <stdbool.h>
#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	char c;
	bool first = true, nl = true;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		nl = c == '\n';
		if (first && c >= 'a' && c <= 'z') {
			putchar(c & 223);
			first = false;
		} else {
			putchar(c);
			first = c == ' ' || c == '\t' || c == '\n';
		}
	}
	if (!nl)
		putchar('\n');
	return 0;
}
