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
		int b, c, d = --a % 26;
		if (a >= 26) {
			a /= 26;
			c = --a % 26;
			if (a >= 26) {
				a /= 26;
				b = --a % 26;
				printf("%c%c%c\n", 'A' + b, 'A' + c, 'A' + d);
			} else {
				printf("%c%c\n", 'A' + c, 'A' + d);
			}
		} else {
			printf("%c\n", 'A' + d);
		}
	}
	return 0;
}
