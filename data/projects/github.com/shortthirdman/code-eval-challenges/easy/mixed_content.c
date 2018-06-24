#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned i, p = 0, pbs = 16, s = 0, sbs = 64;
	char c;
	char *pb = malloc(pbs), *sb = malloc(sbs);
	bool str = false, dig = false, d = true;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || p > 0) {
		if (c == '\n' || c == ',' || c == EOF) {
			if (d) {
				if (s + dig + p >= sbs) {
					do {
						sbs += sbs / 2;
					} while (s + dig + p >= sbs);
					sb = realloc(sb, sbs);
				}
				if (dig)
					sb[s++] = ',';
				else
					dig = true;
				for (i = 0; i < p; i++)
					sb[s++] = pb[i];
			} else {
				if (str)
					printf(",");
				else
					str = true;
				for (i = 0; i < p; i++)
					printf("%c", pb[i]);
			}
			if (c == '\n' || c == EOF) {
				if (dig) {
					if (str) {
						printf("|");
					}
					sb[s] = '\0';
					printf("%s", sb);
				}
				printf("\n");
				s = 0;
				str = false;
				dig = false;
			}
			p = 0;
			d = true;
		} else {
			if (c < '0' || c > '9')
				d = false;
			if (p == pbs) {
				pbs += pbs / 2;
				pb = realloc(pb, pbs);
			}
			pb[p++] = c;
		}
	}
	free(sb);
	free(pb);
	return 0;
}
