#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned s = 0, sbs = 32, p = 0, pbs = 16, m;
	char **pb = malloc(pbs * sizeof(char *));
	char c;
	char *sb = malloc(sbs);
	pb[0] = sb;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || s > 0) {
		if (c == '\n' || c == EOF) {
			sb[s] = '\0';
			sscanf(pb[p], "%u", &m);
			if (p >= m)
				puts(pb[p - m]);
			s = 0;
			p = 0;
		} else {
			if (s == sbs - 1) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
			}
			if (c == ' ') {
				sb[s++] = '\0';
				if (p == pbs - 1) {
					pbs += pbs / 2;
					pb = realloc(pb, pbs * sizeof(char *));
				}
				pb[++p] = sb + s;
			} else {
				sb[s++] = c;
			}
		}
	}
	free(sb);
	free(pb);
	return 0;
}
