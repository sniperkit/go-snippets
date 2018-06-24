#include <stdio.h>

static int happy(int a) {
	int i, ret = 0, t;
	for (i = a; i > 0; i /= 10) {
		t = i % 10;
		ret += t * t;
	}
	return ret;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	int a = 0;
	char c;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || a) {
		if (c == '\n' || c == EOF) {
			int i, j, b[128] = { a, 0 };
			for (i = 0; i < 128 && a > 0; i++) {
				a = happy(a);
				if (a == 1) {
					puts("1");
					break;
				}
				for (j = 0; j < i; j++) {
					if (b[j] == a) {
						a = 0;
						break;
					}
				}
			}
			if (a != 1)
				puts("0");
			a = 0;
		} else {
			a += (c - '0') * (c - '0');
		}
	}
	return 0;
}
