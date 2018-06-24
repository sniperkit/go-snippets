#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int d;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%x", &d) != EOF)
		printf("%d\n", d);
	return 0;
}
