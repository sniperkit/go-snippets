#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	const char *id = "\"id\":";
	char c;
	int i = 0, s = -1, n;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || s >= 0) {
		if (s < 0 && c == '{') {
			s = 0;
		} else if (c == '\n' || c == EOF) {
			if (s >= 0) {
				printf("%d\n", s);
				s = -1;
			}
			i = 0;
		} else if (i == 5) {
			fscanf(fp, " %d", &n);
			if (getc(fp) == ',')
				s += n;
			i = 0;
		} else if (c == id[i]) {
			i++;
		} else {
			i = 0;
		}
	}
	return 0;
}
