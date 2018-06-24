#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int a;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d", &a) != EOF) {
		int l = a, d, n, count = 0, s, t = 0, tx;
		fscanf(fp, "%d %d", &d, &n);
		tx = 6 - d;
		for (s = 6; s <= l - 6; s += d) {
			if (s > tx - d) {
				s = tx;
				if (t == n) {
					tx = l - 6 + d;
				} else {
					fscanf(fp, "%d", &tx);
					t++;
				}
			} else {
				count++;
			}
		}
		printf("%d\n", count);
		while (t++ < n) {
			fscanf(fp, "%d", &a);
		}
	}
	return 0;
}
