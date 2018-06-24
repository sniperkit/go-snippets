#include <stdio.h>

static int nsep(char *p) {
	int i = 0, r = 0;
	while (p[i] != '\0')
		if (p[i++] == ' ')
			r++;
	return r;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	int i = 0;
	char c;
	char line[81];

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || i) {
		if (c == '\n' || c == EOF) {
			line[i] = '\0';
			puts(line);
			i = 0;
		} else {
			line[i] = c;
			if (i == 80) {
				while (i > 0 && line[i] != ' ')
					i--;
				line[i] = '\0';
				if (i == 0) {
					puts(line);
					i = 0;
					line[i++] = c;
				} else {
					int n = nsep(line), j = 0, k = 0, l = 80, h;
					while (line[j] != '\0') {
						if (line[j] != ' ')
							l--;
						j++;
					}
					j = 0;
					while (line[j] != '\0') {
						if (line[j] == ' ') {
							for (h = 0; h < l / n; h++)
								putchar(' ');
							if (k++ < l % n)
								putchar(' ');
						} else {
							putchar(line[j]);
						}
						j++;
					}
					putchar('\n');
					j = i++;
					while (i <= 80) {
						line[i - j - 1] = line[i];
						i++;
					}
					i = 80 - j;
				}
			} else {
				i++;
			}
		}
	}
	return 0;
}
