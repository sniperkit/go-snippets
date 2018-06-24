#include <stdio.h>

static int floors(int e, int s) {
	if (s == 1)
		return 1;
	else if (e == 1)
		return s;
	else if (e == 0 || s == 0)
		return 0;
	else if (e > s)
		return floors(s, s);
	return floors(e - 1, s - 1) + floors(e, s - 1) + 1;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	int e, s, n;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d %d", &e, &s) != EOF) {
		n = 0;
		while (floors(e, n) < s) {
			n++;
		}
		printf("%d\n", n);
	}
	return 0;
}
