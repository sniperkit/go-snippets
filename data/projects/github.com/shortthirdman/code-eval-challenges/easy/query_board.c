#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	int m[256][256] = { {} };
	FILE *fp;
	char b[16];

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fgets(b, 16, fp) != 0) {
		if (b[0] == 'Q') {	/* parameter starts from 9 */
			int i, sum = 0, r = atoi(&b[9]);
			if (b[5] == 'R')	/* row */
				for (i = 0; i < 256; i++)
					sum += m[r][i];
			else			/* column */
				for (i = 0; i < 256; i++)
					sum += m[i][r];
			printf("%d\n", sum);
		} else {	/* first parameter starts from 7 */
			int i, s, r = atoi(&b[7]);
			s = r < 10 ? atoi(&b[9]) : atoi(&b[10]);
			if (b[3] == 'R')	/* row */
				for (i = 0; i < 256; i++)
					m[r][i] = s;
			else			/* column */
				for (i = 0; i < 256; i++)
					m[i][r] = s;
		}
	}
	return 0;
}
