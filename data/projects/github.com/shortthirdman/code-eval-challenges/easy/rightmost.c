#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned sbs = 32, i = 0;
	int j = -1;
	char c, d;
	char *sb = malloc(sbs);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		while (c != ',') {
			if (i == sbs) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
			}
			sb[i++] = c;
			c = getc(fp);
		}
		d = getc(fp);
		while (i > 0)
			if (sb[--i] == d) {
				j = (int)i;
				break;
			}
		printf("%d\n", j);
		fseek(fp, 1, SEEK_CUR);
		i = 0;
		j = -1;
	}
	free(sb);
	return 0;
}
