#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int h;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d ", &h) != EOF) {
		char c;
		int l = 0;
		while (1) {
			int cr = (l + h) % 2;
			fscanf(fp, "%c", &c);
			if (c == 'L') {
				h = (l + h) / 2 + cr - 1;
				fseek(fp, 5, SEEK_CUR);
			} else if (c == 'H') {
				l = (l + h) / 2 + cr + 1;
				fseek(fp, 6, SEEK_CUR);
			} else {
				printf("%d\n", (l + h) / 2 + cr);
				fseek(fp, 4, SEEK_CUR);
				break;
			}
		}
	}
	return 0;
}
