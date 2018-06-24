#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int a;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d", &a) != EOF) {
		int b = 0, c = a;
		while (a) {
			b <<= 1;
			b += a & 1;
			a >>= 1;
		}
		if (c == 0)
			printf("0");
		while (c) {
			printf("%d", b & 1);
			b >>= 1;
			c >>= 1;
		}
		printf("\n");
	}
	return 0;
}
