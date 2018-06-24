#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int a, i, r[20];
	char c;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d ", &r[0]) != EOF) {
		int n = 1;
		while ((c = getc(fp)) != '|') {
			fseek(fp, -1, SEEK_CUR);
			fscanf(fp, "%d ", &r[n++]);
		}
		do {
			fseek(fp, 1, SEEK_CUR);
			for (i = 0; i < n; i++) {
				fscanf(fp, " %d", &a);
				if (a > r[i])
					r[i] = a;
			}
		} while (getc(fp) == ' ');
		printf("%d", r[0]);
		for (i = 1; i < n; i++)
			printf(" %d", r[i]);
		putchar('\n');
	}
	return 0;
}
