#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned i = 0, j, sbs = 32;
	char c;
	char *sb = malloc(sbs);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		if (c == ',') {
			fseek(fp, 1, SEEK_CUR);
			while ((c = getc(fp)) != '\n' && c != EOF)
				for (j = 0; j < i; j++)
					if (sb[j] == c)
						sb[j] = ',';
			for (j = 0; j < i; j++)
				if (sb[j] != ',')
					printf("%c", sb[j]);
			printf("\n");
			i = 0;
		} else {
			if (i == sbs - 1) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
			}
			sb[i++] = c;
		}
	}
	free(sb);
	return 0;
}
