#include <stdbool.h>
#include <stdio.h>

static bool pali(int a) {
	int l[16], c = a, x = 0, y = 0;
	while (c) {
		l[y++] = c % 10;
		c /= 10;
	}
	while (l[x] == l[y - 1]) {
		if (y - x < 3)
			return true;
		x++;
		y--;
	}
	return false;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	int l, r;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d %d", &l, &r) != EOF) {
		int i, n = 0;
		for (i = l; i <= r; i++) {
			int j, prev = -1;
			for (j = i; j <= r; j++) {
				int k, p;
				if (prev > -1) {
					p = prev;
					if (pali(j))
						p++;
				} else {
					p = 0;
					for (k = i; k <= j; k++) {
						if (pali(k))
							p++;
					}
				}
				if (p % 2 == 0)
					n++;
				prev = p;
			}
		}
		printf("%d\n", n);
	}
	return 0;
}
