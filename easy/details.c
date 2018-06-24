#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	char c, p = '\n';
	int m = 10, dots = 0;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || p != '\n') {
		switch (c) {
		case '\n':
		case EOF:
			if (dots < m)
				m = dots;
			printf("%d\n", m);
			p = '\n';
			m = 10;
			dots = 0;
			break;
		case '.':
			if (p != 'Y')
				dots++;
			break;
		case ',':
			if (dots < m)
				m = dots;
			dots = 0;
		default:
			p = c;
		}
	}
	return 0;
}
