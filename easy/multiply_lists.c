#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	char c;
	int b[10] = { 0 };

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%c", &c) != EOF) {
		int a = 0, n = 0, i;
		do {
			if (c < '0' || c > '9') {
				goto fail_nondigit;
			}
			a = 0;
			while (c >= '0' && c <= '9') {
				a = a * 10 + c - '0';
				if (fscanf(fp, "%c", &c) == EOF)
					goto fail_eof;
			}
			b[n++] = a;
			if (fscanf(fp, "%c", &c) == EOF)
				goto fail_eof;
		} while (c != '|');
		if (fscanf(fp, "%c", &c) == EOF)
			goto fail_eof;
		if (c != ' ')
			goto fail_nonspace;
		if (fscanf(fp, "%c", &c) == EOF)
			goto fail_eof;
		for (i = 0; i < n; i++) {
			if (c < '0' || c > '9') {
				goto fail_nondigit;
			}
			a = 0;
			while (c >= '0' && c <= '9') {
				a = a * 10 + c - '0';
				if (fscanf(fp, "%c", &c) == EOF) {
					if (i < n - 1)
						goto fail_eof;
					else
						break;
				}
			}
			if (i > 0)
				printf(" ");
			printf("%d", b[i] * a);
			if (i < n - 1 && fscanf(fp, "%c", &c) == EOF)
				goto fail_eof;
		}
		printf("\n");
	}
	return 0;
fail_nondigit:
	printf("unexpected non-digit input: \"%c\"\n", c);
	return EXIT_FAILURE;
fail_nonspace:
	printf("unexpected non-space input: \"%c\"\n", c);
	return EXIT_FAILURE;
fail_eof:
	printf("unexpected end of input\n");
	return EXIT_FAILURE;
}
