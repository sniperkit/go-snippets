#include <stdio.h>

static int sq(int a) {
	int i;
	for (i = 0; (i + 1) * (i + 1) <= a; i++) ;
	return i;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	int a[4], x, y, d;
	char line[27];

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fgets(line, 27, fp) != 0) {
		sscanf(line, "(%d, %d) (%d, %d)", &a[0], &a[1], &a[2], &a[3]);
		x = a[0] - a[2];
		y = a[1] - a[3];
		d = sq(x * x + y * y);
		printf("%d\n", d);
	}
	return 0;
}
