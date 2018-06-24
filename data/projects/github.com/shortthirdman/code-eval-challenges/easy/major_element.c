#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	char c;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		int *z = calloc(101, sizeof(int));
		int n = 0, mn = 0, mc = 0;
		do {
			int num = 0;
			n++;
			while (c >= '0' && c <= '9') {
				num *= 10;
				num += c - '0';
				c = getc(fp);
			}
			z[num]++;
			if (z[num] > mc) {
				mc++;
				mn = num;
			}
			if (c == ',')
				c = getc(fp);
		} while (c >= '0' && c <= '9');
		free(z);
		if (mc > n / 2)
			printf("%d\n", mn);
		else
			puts("None");
	}
	return 0;
}
