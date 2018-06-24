#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned i, s = 0, sbs = 16;
	char c;
	float f, g;
	char *sb = malloc(sbs);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		while (c == '*' || c == '/' || c == '+') {
			if (s == sbs) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
			}
			sb[s++] = c;
			fseek(fp, 1, SEEK_CUR);
			c = getc(fp);
		}
		fseek(fp, -1, SEEK_CUR);
		fscanf(fp, "%f ", &f);
		for (i = s; i > 0; i--) {
			fscanf(fp, "%f ", &g);
			switch (sb[i - 1]) {
			case '*':
				f *= g;
				break;
			case '/':
				f /= g;
				break;
			default:
				f += g;
			}
		}
		printf("%d\n", (int)(f + 0.001));
		s = 0;
	}

	free(sb);
	return 0;
}
