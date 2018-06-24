#include <stdio.h>

static int rev(int a) {
	int r = 0, s;
	for (s = a; s > 0; s /= 10)
		r = 10 * r + s % 10;
	return r;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	int a, i, r;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d", &a) != EOF) {
		for (i = 0; i < 100; i++) {
			r = rev(a);
			if (r == a)
				break;
			a += r;
		}
		if (i < 100)
			printf("%d %d\n", i, a);
		else
			puts("not found");
	}
	return 0;
}
