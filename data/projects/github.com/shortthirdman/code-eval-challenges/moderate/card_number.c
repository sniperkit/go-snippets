#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int i = 0, e = 0, o = 0;
	char c;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || i) {
		if (c == '\n' || c == EOF) {
			if (i % 2) {
				printf("%d\n", o % 10 == 0);
			} else {
				printf("%d\n", e % 10 == 0);
			}
			i = 0;
			e = 0;
			o = 0;
		} else if (c != ' ') {
			int x = c - '0';
			if (i % 2) {
				e += x;
				x *= 2;
				if (x > 9) {
					o += x % 10 + 1;
				} else {
					o += x;
				}
			} else {
				o += x;
				x *= 2;
				if (x > 9) {
					e += x % 10 + 1;
				} else {
					e += x;
				}
			}
			i++;
		}
	}
	return 0;
}
