#include <stdio.h>
#include <stdlib.h>

static int max(int a, int b) {
	return a > b ? a : b;
}

static char *bt(int *ib, char *sb, int s, int t, int i, int j, int p) {
	if (i < 1 || j < 1)
		return sb + p;
	if (sb[i - 1] == sb[s + j - 1]) {
		sb[--p] = sb[i - 1];
		return bt(ib, sb, s, t, i - 1, j - 1, p);
	}
	if (ib[i * (t + 1) + j - 1] > ib[(i - 1) * (t + 1) + j])
		return bt(ib, sb, s, t, i, j - 1, p);
	return bt(ib, sb, s, t, i - 1, j, p);
}

int main(int argc, char *argv[]) {
	FILE *fp;
	int s = 0, t = -1, sbs = 45, i, j, ibs = 514;
	int *ib = malloc(ibs * sizeof(int));
	char c;
	char *sb = malloc(sbs + 1);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || s > 0) {
		if (c == '\n' || c == EOF) {
			if (t < 0) {
				s = 0;
				continue;
			}
			i = (s + 1) * (t + 1);
			if (i > ibs) {
				ibs = i;
				ib = realloc(ib, ibs * sizeof(int));
			}
			while (--i >= 0)
				ib[i] = 0;
			for (i = 0; i < s; i++) {
				for (j = 0; j < t; j++)
					if (sb[i] == sb[s + j])
						ib[(i + 1) * (t + 1) + j + 1] = ib[i * (t + 1) + j] + 1;
					else
						ib[(i + 1) * (t + 1) + j + 1] = max(ib[(i + 1) * (t + 1) + j],
										    ib[i * (t + 1) + j + 1]);
			}
			sb[s + t] = '\0';
			puts(bt(ib, sb, s, t, s, t, s + t));
			s = 0;
			t = -1;
		} else if (c == ';') {
			t = 0;
		} else if (t < 0) {
			if (s == sbs) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs + 1);
			}
			sb[s++] = c;
		} else {
			if (s + t == sbs) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs + 1);
			}
			sb[s + t++] = c;
		}
	}
	free(sb);
	free(ib);
	return 0;
}
