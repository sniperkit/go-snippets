#include <stdio.h>

static int sq(int a) {
	int i;
	for (i = 1; i * i < a; i++) ;
	return i;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	char c, m[100];
	int n = 0, l, i, j;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || n) {
		if (c == '\n' || c == EOF) {
			l = sq(n);
			for (i = 0; i < l; i++)
				for (j = l - 1; j >= 0; j--) {
					putchar(m[j * l + i]);
					if (i != l - 1 || j > 0)
						putchar(' ');
					else
						putchar('\n');
				}
			n = 0;
		} else if (c != ' ') {
			m[n++] = c;
		}
	}
	return 0;
}
