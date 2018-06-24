#include <stdbool.h>
#include <stdio.h>

static bool isMagic(int a) {
	int d = 0, i, r, l = 0, ns[4];
	while (a) {
		r = a % 10;
		if (r == 0 || (d & (1 << r)) > 0)
			return false;
		d |= 1 << r;
		ns[l++] = r;
		a /= 10;
	}
	d = 0;
	r = 0;
	for (i = 0; i < l; i++) {
		r = (r + ns[l - 1 - r]) % l;
		if (d & (1 << r))
			return false;
		d |= 1 << r;
	}
	return r == 0;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	int t = 1, a, b, i, ibs = 89;

	int ib[89];
	for (i = 0; i < ibs; i++) {
		while (!isMagic(t))
			t++;
		ib[i] = t++;
	}

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d %d", &a, &b) != EOF) {
		bool first = true;
		for (i = 0; i < ibs && ib[i] <= b; i++) {
			if (ib[i] >= a) {
				if (!first)
					putchar(' ');
				else
					first = false;
				printf("%d", ib[i]);
			}
		}
		printf(first ? "-1\n" : "\n");
	}
	return 0;
}
