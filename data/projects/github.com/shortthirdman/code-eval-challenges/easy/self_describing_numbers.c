#include <stdbool.h>
#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int p = 0, a[10] = { 0 }, b[10] = { 0 };
	char c;
	bool s = true;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || p > 0) {
		if (c == '\n' || c == EOF) {
			for (p = 9; p >= 0; p--) {
				if (a[p] != b[p])
					s = false;
				a[p] = 0;
				b[p] = 0;
			}
			if (s)
				puts("1");
			else
				puts("0");
			p = 0;
			s = true;
		} else {
			int v = c - '0';
			if (p > 9) {
				if (v > 0)
					s = false;
				b[0]++;
			} else {
				a[p] = v;
				b[a[p++]]++;
			}
		}
	}
	return 0;
}
