#include <math.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	char line[13];

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fgets(line, 7, fp) != 0) {
		int i, j, n = atoi(line);
		unsigned long clos2 = 100000000;
		if (n == 0)
			return 0;
		int a[n][2];

		for (i = 0; i < n; i++) {
			fgets(line, 13, fp);
			sscanf(line, "%d %d", &a[i][0], &a[i][1]);
		}

		if (n > 1) {
			for (i = 0; i < n - 1; i++) {
				for (j = i + 1; j < n; j++) {
					int x = a[i][0] - a[j][0];
					int y = a[i][1] - a[j][1];
					unsigned long dist2 = x * x + y * y;
					if (dist2 < clos2)
						clos2 = dist2;
				}
			}
		}

		if (clos2 == 100000000)
			printf("INFINITY\n");
		else
			printf("%.4f\n", sqrtf((float)clos2));
	}
	return 0;
}
