#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int p = -1, cp = -1, rd = 0, i;
	char c, d = '|', sb[13];
	sb[12] = '\0';

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		for (i = 0; i < 12; i++) {
			if (p == -1 || (p - i < 2 && p - i > -2)) {
				if (c == 'C')
					cp = i;
				else if (c == '_')
					rd = i;
			}
			sb[i] = c;
			c = getc(fp);
		}
		if (p == -1) {
			if (cp != -1)
				p = cp;
			else
				p = rd;
		} else {
			if (cp != -1) {
				if (cp < p)
					d = '/';
				else if (cp > p)
					d = '\\';
				else
					d = '|';
				p = cp;
			} else {
				if (rd < p)
					d = '/';
				else if (rd > p)
					d = '\\';
				else
					d = '|';
				p = rd;
			}
		}
		sb[p] = d;
		puts(sb);
		cp = -1;
	}
	return 0;
}
