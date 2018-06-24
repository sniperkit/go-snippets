#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int i = 0, l = 40;
	char c, s[55];

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		if (c == '\n') {
			s[i] = '\0';
			puts(s);
			i = 0;
			l = 40;
		} else if (i == 55) {
			s[l] = '\0';
			printf("%s", s);
			puts("... <Read More>");
			i = 0;
			l = 40;
			while (c != '\n' && c != EOF)
				c = getc(fp);
		} else {
			if (i < 40 && c == ' ')
				l = i;
			s[i++] = c;
		}
	}
	return 0;
}
