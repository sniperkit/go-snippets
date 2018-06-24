#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	float angle;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%f", &angle) != EOF) {
		printf("%d.", (int)angle);
		angle = (angle - (int)angle) * 60;
		printf("%02d'", (int)angle);
		angle = (angle - (int)angle) * 60;
		printf("%02d\"\n", (int)angle);
	}
	return 0;
}
