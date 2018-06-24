#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	char c;
	unsigned i, j, k, r, t, ll[32], cl[32];

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		r = 0;
		for (j = 0; j < 32; j++) {
			ll[j] = 0;
			cl[j] = 0;
		}
		i = 0;
		j = 1;
		do {
			if (c == '|') {
				for (k = 1; k < 31; k++) {
					ll[k] = cl[k];
					cl[k] = 0;
				}
				i++;
				j = 0;
			} else if (c == 'o') {
				if (ll[j - 1] > 0)
					cl[j] = ll[j - 1];
				else if (ll[j] > 0)
					cl[j] = ll[j];
				else if (cl[j - 1] > 0)
					cl[j] = cl[j - 1];
				if (ll[j + 1] > 0) {
					if (cl[j] == 0) {
						cl[j] = ll[j + 1];
					} else if (cl[j] != ll[j + 1]) {
						t = ll[j + 1];
						for (k = j + 1; k < 31; k++)
							if (ll[k] == t)
								ll[k] = cl[j];
						for (k = 1; k < j - 1; k++)
							if (cl[k] == t)
								cl[k] = cl[j];
						r--;
					}
				}
				if (cl[j] == 0) {
					cl[j] = i * 32 + j;
					r++;
				}
			}
			j++;
			if ((c = getc(fp)) == EOF || c == '\n')
				break;
			c = getc(fp);
		} while (c != '\n' && c != EOF);
		printf("%d\n", r);
	}
	return 0;
}
