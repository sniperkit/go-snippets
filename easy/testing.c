#include <stdio.h>

static const char *m[] = { "Done", "Low", "Medium", "High", "Critical" };

static int prio(int a) {
	switch (a) {
	case 0:
		return 0;
	case 1:
	case 2:
		return 1;
	case 3:
	case 4:
		return 2;
	case 5:
	case 6:
		return 3;
	default:
		return 4;
	}
}

int main(int argc, char *argv[]) {
	FILE *fp;
	char c, d[40];
	int s = 0, i = 0, r = 0;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || s > 0) {
		switch (s) {
		case 0:
			if (c == ' ')
				s = 1;
			else
				d[i++] = c;
			break;
		case 1:
			if (c != '|') {
				d[i++] = ' ';
				d[i++] = c;
				s = 0;
			} else {
				fseek(fp, 1, SEEK_CUR);
				s = 2;
				i = 0;
			}
			break;
		case 2:
			if (c == '\n' || c == EOF) {
				puts(m[prio(r)]);
				s = 0;
				i = 0;
				r = 0;
			} else if (d[i++] != c) {
				r++;
			}
		}
	}
	return 0;
}
