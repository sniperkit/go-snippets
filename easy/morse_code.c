#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	char c;
	int m = 1;
	const char *morse = "ETIANMSURWDKGOHVF L PJBXCYZQ  54 3   2       16       7   8 90";

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || m > 1) {
		if (c == '.') {
			m <<= 1;
		} else if (c == '-') {
			m = (m << 1) + 1;
		} else if (m == 1) {
			putchar(' ');
		} else {
			if (m < 64)
				putchar(morse[m - 2]);
			if (c == '\n' || c == EOF)
				putchar('\n');
			m = 1;
		}
	}
	return 0;
}
