#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int s = 0, a;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d", &a) != EOF)
		s += a;
	printf("%d\n", s);
	return 0;
}
