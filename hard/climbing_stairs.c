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
		if (a == 0 || a == 1) {
			printf("%d\n", a);
		} else {
			int b[3] = { 1, 1, 0 }, c = 0;
			while (a > c++) {
				b[0] = b[1] + b[2];
				b[2] = b[1];
				b[1] = b[0];
			}
			printf("%d\n", b[0]);
		}
	}
	return 0;
}
