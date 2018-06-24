#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned d;
	int s = 0;
	char c;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%x", &d) != EOF) {
		s += d;
		while (fscanf(fp, "%x", &d) > 0)
			s += d;
		fseek(fp, 2, SEEK_CUR);
		d = 0;
		do {
			c = getc(fp);
			if (c == '0' || c == '1') {
				d = 2 * d + c - '0';
			} else {
				s -= d;
				d = 0;
			}
		} while (c != '\n' && c != EOF);
		puts(s > 0 ? "False" : "True");
		s = 0;
	}
	return 0;
}
