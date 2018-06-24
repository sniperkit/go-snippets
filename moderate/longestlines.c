#include <stdio.h>
#include <stdlib.h>

struct lbuffer {
	unsigned length;
	char *buffer;
};

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned n, i, s = 0, sbs = 32, t;
	char c;
	char *sb = malloc(sbs), *tmp;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	fscanf(fp, "%u\n", &n);
	struct lbuffer *lb = malloc(n * sizeof(struct lbuffer));
	for (i = 0; i < n; i++) {
		lb[i].length = 0;
		lb[i].buffer = malloc(sbs);
	}
	while ((c = getc(fp)) != EOF || s > 0) {
		if (c == '\n' || c == EOF) {
			sb[s] = '\0';
			for (i = 0; i < n; i++) {
				if (s > lb[i].length) {
					tmp = sb;
					t = s;
					sb = lb[i].buffer;
					s = lb[i].length;
					lb[i].buffer = tmp;
					lb[i].length = t;
				}
			}
			s = 0;
		} else {
			if (s == sbs - 1) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
				for (i = 0; i < n; i++)
					lb[i].buffer = realloc(lb[i].buffer, sbs);
			}
			sb[s++] = c;
		}
	}
	free(sb);
	for (i = 0; i < n; i++) {
		puts(lb[i].buffer);
		free(lb[i].buffer);
	}
	return 0;
}
