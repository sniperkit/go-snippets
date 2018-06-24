#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int n, m;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d,%d", &n, &m) != EOF) {
		printf("%d\n", n - (n / m) * m);
	}
	return 0;
}
