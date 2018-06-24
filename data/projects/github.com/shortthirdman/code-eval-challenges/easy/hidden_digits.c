#include <stdbool.h>
#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	char c;
	bool n = false;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || n) {
		if (c == '\n' || c == EOF) {
			if (!n)
				printf("NONE");
			printf("\n");
			n = false;
		} else if (c >= '0' && c <= '9') {
			printf("%c", c);
			n = true;
		} else if (c >= 'a' && c <= 'j') {
			printf("%c", c - 'a' + '0');
			n = true;
		}
	}
	return 0;
}
