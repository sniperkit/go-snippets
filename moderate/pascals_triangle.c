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
		int i;
		printf("1");
		for (i = 1; i < a; i++) {
			int j, r = 1;
			printf(" 1");
			for (j = 1; j <= i; j++) {
				r = (r * (i + 1 - j)) / j;
				printf(" %d", r);
			}
		}
		printf("\n");
	}
	return 0;
}
