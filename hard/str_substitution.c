#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned ss = 0, sbs = 192, cs1 = 0, cs2 = 0, cb1s = 32, cb2s = 32, i, j, d;
	char *sb = malloc(sbs), *cb1 = malloc(cb1s), *cb2 = malloc(cb2s);
	char c;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		while (c != ';') {
			if (ss == sbs - 1) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
			}
			sb[ss++] = c;
			c = getc(fp);
		}

		do {
			while ((c = getc(fp)) != ',') {
				if (cs1 == cb1s) {
					cb1s += cb1s / 2;
					cb1 = realloc(cb1, cb1s);
				}
				cb1[cs1++] = c;
			}
			while ((c = getc(fp)) == '0' || c == '1') {
				if (cs2 == cb2s) {
					cb2s += cb2s / 2;
					cb2 = realloc(cb2, cb2s);
				}
				cb2[cs2++] = c - '0' + 'a';
			}
			if (cs2 == 0) {
				cb2[0] = 'x';
				cs2 = 1;
			}
			for (i = 0; i < ss - cs1 + 1; i++) {
				bool f = true;
				for (j = 0; j < cs1; j++) {
					if (sb[i + j] != cb1[j]) {
						f = false;
						break;
					}
				}
				if (!f) {
					continue;
				}
				if (cs1 < cs2) {
					d = cs2 - cs1;
					while (ss + d >= sbs - 1) {
						sbs += sbs / 2;
						sb = realloc(sb, sbs);
					}
					ss += d;
					for (j = ss - 1; j >= i + cs1 + d; j--) {
						sb[j] = sb[j - d];
					}
				} else if (cs1 > cs2) {
					d = cs1 - cs2;
					for (j = i + cs2; j < ss - d; j++) {
						sb[j] = sb[j + d];
					}
					ss -= d;
				}
				for (j = 0; j < cs2; j++) {
					sb[i + j] = cb2[j];
				}
			}

			cs1 = 0;
			cs2 = 0;
		} while (c == ',');
		sb[ss] = '\0';
		for (i = ss; i > 0; i--) {
			switch (sb[i - 1]) {
			case 'a':
				sb[i - 1] = '0';
				break;
			case 'b':
				sb[i - 1] = '1';
				break;
			case 'x':
				for (j = i - 1; j < ss; j++) {
					sb[j] = sb[j + 1];
				}
				ss--;
			}
		}
		printf("%s\n", sb);
		ss = 0;
	}
	free(sb);
	free(cb1);
	free(cb2);
	return 0;
}
