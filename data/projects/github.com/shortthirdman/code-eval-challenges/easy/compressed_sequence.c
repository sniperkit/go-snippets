#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int a, b = -1, c = 1;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d", &a) != EOF) {
		if (a == b) {
			c++;
		} else {
			if (b >= 0) {
				printf("%d %d ", c, b);
				c = 1;
			}
			b = a;
		}
		if (getc(fp) != ' ') {
			printf("%d %d\n", c, b);
			b = -1;
			c = 1;
		}
	}
	return 0;
}
