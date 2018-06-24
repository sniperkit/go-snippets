#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned sbs = 32, s = 0, ibs = 32, i = 0, j, n, m;
	char c;
	char *sb = malloc(sbs);
	unsigned *ib = malloc(ibs * sizeof(unsigned));
	ib[0] = 0;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((fscanf(fp, "%d;%d;", &n, &m)) != EOF) {
		unsigned tn = 0, te, ts, tw = 0;
		if (n * m > ibs) {
			ibs = n * m;
			ib = realloc(ib, ibs * sizeof(unsigned));
		}
		while (i < n * m) {
			if (s == sbs) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
			}
			if ((c = getc(fp)) == ' ' || c == '\n' || c == EOF) {
				c = '\0';
				if (++i < n * m)
					ib[i] = s + 1;
			}
			sb[s++] = c;
		}
		i = 0;
		j = 1;
		te = m - 1;
		ts = n - 1;
		printf("%s", sb);
		do {
			while (j <= te) {
				printf(" %s", sb + ib[i * m + j]);
				if (j < te)
					j++;
				else
					break;
			}
			i++;
			tn++;
			if (i > ts)
				break;
			while (i <= ts) {
				printf(" %s", sb + ib[i * m + j]);
				if (i < ts)
					i++;
				else
					break;
			}
			j--;
			te--;
			if (j < tw)
				break;
			while (j >= tw) {
				printf(" %s", sb + ib[i * m + j]);
				if (j > tw)
					j--;
				else
					break;
			}
			i--;
			ts--;
			if (i < tn)
				break;
			while (i >= tn) {
				printf(" %s", sb + ib[i * m + j]);
				if (i > tn)
					i--;
				else
					break;
			}
			j++;
			tw++;
		} while (j <= te);
		printf("\n");
		s = 0;
		i = 0;
	}
	free(ib);
	free(sb);
	return 0;
}
