#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int a;
	unsigned b, c;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d,%d,%d", &a, &b, &c) != EOF)
		if (!!(a & 1 << --b) == !!(a & 1 << --c))
			puts("true");
		else
			puts("false");
	return 0;
}
