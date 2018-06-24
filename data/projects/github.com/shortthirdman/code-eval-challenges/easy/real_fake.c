#include <stdio.h>

static const char *m[] = { "Fake", "Real" };

int main(int argc, char *argv[]) {
	FILE *fp;
	int i = 0, e = 0;
	char c;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || i) {
		if (c == '\n' || c == EOF) {
			printf("%s\n", m[e % 10 == 0]);
			i = 0;
			e = 0;
		} else if (c != ' ') {
			int x = c - '0';
			if (i & 1)
				e += x;
			else
				e += x * 2;
			i++;
		}
	}
	return 0;
}
