#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	const int stats[32] = { 2, 4, 8, 6, 3, 9, 7, 1, 4, 6, 4, 6, 5, 5, 5, 5,
				6, 6, 6, 6, 7, 9, 3, 1, 8, 4, 2, 6, 9, 1, 9, 1 };
	int a;
	long long int n;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d %lld", &a, &n) != EOF) {
		int i;
		long long int m = n / 4;
		long long int res[10] = { 0 };
		a -= 2;
		for (i = 0; i < 4; i++) {
			res[stats[a * 4 + i]] += m;
		}
		for (i = 0; i < n % 4; i++) {
			res[stats[a * 4 + i]]++;
		}
		for (i = 0; i <= 9; i++) {
			if (i > 0) {
				printf(", ");
			}
			printf("%d: %lld", i, res[i]);
		}
		printf("\n");
	}
	return 0;
}
