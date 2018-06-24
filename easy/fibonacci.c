#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned a, t;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%u", &a) != EOF) {
		if (a < 2) {
			printf("%d\n", a);
		} else {
			unsigned b = 0, c = 1;
			for (; a > 1; a--) {
				t = b + c;
				b = c;
				c = t;
			}
			printf("%u\n", c);
		}
	}
	return 0;
}
