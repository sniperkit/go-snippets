#include <stdio.h>

static const int m[5] = { 50, 25, 10, 5, 1 };

static int alter(int n, int p) {
	if (n == 0)
		return 1;
	while (m[p] > n)
		p++;
	if (m[p] == 1)
		return 1;
	return alter(n - m[p], p) + alter(n, p + 1);
}

int main(int argc, char *argv[]) {
	FILE *fp;
	int n;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d;", &n) != EOF)
		printf("%d\n", alter(n, 0));
	return 0;
}
