#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	char c;
	float count = 0, lower = 0, ratio;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || count > 0) {
		if (c == '\n' || c == EOF) {
			if (count > 0) {
				ratio = 100 * lower / count;
				printf("lowercase: %.2f uppercase: %.2f\n",
				       ratio, 100 - ratio);
				count = 0;
				lower = 0;
			} else {
				puts("lowercase: 0.00 uppercase: 0.00");
			}
		} else if (c >= 'a' && c <= 'z') {
			lower++;
			count++;
		} else if (c >= 'A' && c <= 'Z') {
			count++;
		}
	}
	return 0;
}
