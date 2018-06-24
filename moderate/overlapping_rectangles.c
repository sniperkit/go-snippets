#include <stdbool.h>
#include <stdio.h>

static bool within(int *a, int r, int p, int q) {
	return a[p] >= a[r] && a[p] <= a[r + 2] &&
	       a[q] >= a[r + 3] && a[q] <= a[r + 1];
}

int main(int argc, char *argv[]) {
	FILE *fp;
	int a[8];

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d,%d,%d,%d,%d,%d,%d,%d\n", &a[0], &a[1],
		      &a[2], &a[3], &a[4], &a[5], &a[6], &a[7]) != EOF) {
		if (within(a, 0, 4, 5) || within(a, 0, 6, 7) ||
		    within(a, 0, 4, 7) || within(a, 0, 6, 5) ||
		    within(a, 4, 0, 1) || within(a, 4, 2, 3) ||
		    within(a, 4, 0, 3) || within(a, 4, 2, 1))
			puts("True");
		else
			puts("False");
	}
	return 0;
}
