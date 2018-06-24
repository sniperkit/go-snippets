#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned a;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%u", &a) != EOF)
		printf("%u\n", __builtin_popcount(a) % 3);
	return 0;
}
