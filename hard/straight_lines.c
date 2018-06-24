#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	char c;
	unsigned ps = 32;
	int *p = malloc(ps * 2 * sizeof(int)), dx, dy;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d %d", &p[0], &p[1]) != EOF) {
		unsigned ip = 1, count = 0, i, j, k;
		while ((c = getc(fp)) == ' ') {
			if (ip == ps) {
				ps += ps / 2;
				p = realloc(p, ps * 2 * sizeof(int));
			}
			fscanf(fp, "| %d %d", &p[2 * ip], &p[2 * ip + 1]);
			ip++;
		}
		for (i = 0; i < ip - 2; i++)
			for (j = i + 1; j < ip - 1; j++) {
				dx = p[i * 2] - p[j * 2];
				dy = p[i * 2 + 1] - p[j * 2 + 1];
				for (k = 0; k < ip; k++)
					if (k != i && k != j &&
					    dx * (p[i * 2 + 1] - p[k * 2 + 1]) == (p[i * 2] - p[k * 2]) * dy) {
						if (k > j)
							count++;
						break;
					}
			}
		printf("%d\n", count);
	}
	free(p);
	return 0;
}
