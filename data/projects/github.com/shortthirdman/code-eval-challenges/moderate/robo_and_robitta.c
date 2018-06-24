#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned a, b, x, y;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%dx%d | %d %d", &a, &b, &x, &y) != EOF) {
		unsigned r = 0, t;
		while (b > y) {
			r += a;
			t = a;
			a = b - 1;
			b = t;
			t = x;
			x = a + 1 - y;
			y = t;
		}
		printf("%d\n", r + x);
	}
	return 0;
}
