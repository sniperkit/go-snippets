#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	char c;
	int num = 0, ar = 0, rr = 0;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || rr > 0) {
		int a, r = 0;
		if (c == '\n' || c == EOF) {
			printf("%d\n", num + ar * rr);
			num = 0;
			ar = 0;
			rr = 0;
			continue;
		}
		a = c - '0';
		switch (getc(fp)) {
		case 'M':
			r = 1000;
			break;
		case 'D':
			r = 500;
			break;
		case 'C':
			r = 100;
			break;
		case 'L':
			r = 50;
			break;
		case 'X':
			r = 10;
			break;
		case 'V':
			r = 5;
			break;
		case 'I':
			r = 1;
		}
		if (r > rr) {
			num -= ar * rr;
		} else {
			num += ar * rr;
		}
		ar = a;
		rr = r;
	}
	return 0;
}
