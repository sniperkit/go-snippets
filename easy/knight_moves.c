#include <stdbool.h>
#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	char a, b;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%c%c\n", &a, &b) != EOF) {
		int l = a - 'a', r = b - '1';
		bool first = true;
		if (l > 1 && r > 0) {
			printf("%c%c", 'a' + l - 2, '1' + r - 1);
			first = false;
		}
		if (l > 1 && r < 7) {
			if (!first)
				putchar(' ');
			printf("%c%c", 'a' + l - 2, '1' + r + 1);
			first = false;
		}
		if (l > 0 && r > 1) {
			if (!first)
				putchar(' ');
			printf("%c%c", 'a' + l - 1, '1' + r - 2);
			first = false;
		}
		if (l > 0 && r < 6) {
			if (!first)
				putchar(' ');
			printf("%c%c", 'a' + l - 1, '1' + r + 2);
			first = false;
		}
		if (l < 7 && r > 1) {
			if (!first)
				putchar(' ');
			printf("%c%c", 'a' + l + 1, '1' + r - 2);
			first = false;
		}
		if (l < 7 && r < 6) {
			if (!first)
				putchar(' ');
			printf("%c%c", 'a' + l + 1, '1' + r + 2);
		}
		if (l < 6 && r > 0)
			printf(" %c%c", 'a' + l + 2, '1' + r - 1);
		if (l < 6 && r < 7)
			printf(" %c%c", 'a' + l + 2, '1' + r + 1);
		putchar('\n');
	}
	return 0;
}
