#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int a, b, i = 0;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d,%d", &a, &b) != EOF) {
		while (b * ++i < a);
		printf("%d\n", b * i);
		i = 0;
	}
	return 0;
}
