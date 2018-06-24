#include <stdbool.h>
#include <stdio.h>

static const char *units[] = {
	"PENNY", "NICKEL", "DIME", "QUARTER", "HALF DOLLAR", "ONE",
	"TWO", "FIVE", "TEN", "TWENTY", "FIFTY", "ONE HUNDRED"
};
static const int value[] = {
	1, 5, 10, 25, 50, 100,
	200, 500, 1000, 2000, 5000, 10000
};
static const int sorted[] = {
	2, 10, 7, 4, 1, 5, 11, 0, 3, 8, 9, 6
};

int main(int argc, char *argv[]) {
	FILE *fp;
	int p, p1, c, c1, i;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d", &p1) != EOF) {
		int p2 = 0, c2 = 0;
		if (getc(fp) == '.')
			fscanf(fp, "%d%*c", &p2);
		fscanf(fp, "%d", &c1);
		if (getc(fp) == '.')
			fscanf(fp, "%d%*c", &c2);
		p = p1 * 100 + p2;
		c = c1 * 100 + c2;
		if (p > c) {
			printf("ERROR");
		} else if (p == c) {
			printf("ZERO");
		} else {
			int change[12] = { 0 };
			bool con = false;
			for (i = 11; i >= 0; i--) {
				while (c - p >= value[i]) {
					change[i]++;
					c -= value[i];
				}
			}
			for (i = 0; i < 12; i++) {
				while (change[sorted[i]]) {
					if (con) {
						printf(",");
					} else {
						con = true;
					}
					printf("%s", units[sorted[i]]);
					change[sorted[i]]--;
				}
			}
		}
		printf("\n");
	}
	return 0;
}
