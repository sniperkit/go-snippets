#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned s = 0, sbs = 32, n = 0, nbs = 16, t, i;
	char c;
	char *sb = malloc(sbs);
	int *nb = malloc(nbs * sizeof(int));
	int *pb = malloc(nbs * sizeof(int));
	nb[0] = 0;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		switch (c) {
		case ';':
			sb[s] = '\0';
			for (i = 0; i <= n; i++)
				pb[i] = -1;
			for (i = 0; i < n; i++) {
				fscanf(fp, "%d", &t);
				pb[t - 1] = nb[i];
				nb[i] = -1;
			}
			for (i = 0; i <= n; i++) {
				if (pb[i] == -1)
					pb[i] = nb[n];
				if (i > 0)
					printf(" ");
				printf("%s", sb + pb[i]);
			}
			printf("\n");
			s = 0;
			n = 0;
			nb[0] = 0;
			fseek(fp, 1, SEEK_CUR);
			break;
		case ' ':
			if (s == sbs - 1) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
			}
			sb[s++] = '\0';
			if (n == nbs - 1) {
				nbs += nbs / 2;
				nb = realloc(nb, nbs * sizeof(int));
				pb = realloc(pb, nbs * sizeof(int));
			}
			nb[++n] = (int)s;
			break;
		default:
			if (s == sbs - 1) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
			}
			sb[s++] = c;
		}
	}
	free(pb);
	free(nb);
	free(sb);
	return 0;
}
