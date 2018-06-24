#include <stdbool.h>
#include <stdio.h>

unsigned g2d(unsigned a) {
	a ^= a >> 4;
	a ^= a >> 2;
	a ^= a >> 1;
	return a;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	bool p = true;
	char c;
	unsigned n = 0;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || !p) {
		p = false;
		if (c == '0' || c == '1') {
			n = 2 * n + c - '0';
		} else if (c == ' ') {
			printf("%u | ", g2d(n));
			n = 0;
			fseek(fp, 2, SEEK_CUR);
		} else if (c == '\n' || c == EOF) {
			printf("%u\n", g2d(n));
			n = 0;
			p = true;
		}
	}
	return 0;
}
