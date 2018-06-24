#include <stdio.h>

static const char *v[] = { "2", "3", "4", "5", "6", "7", "8", "9", "10",
			   "J", "Q", "K", "A" };

static unsigned rank(char c) {
	unsigned i;
	for (i = 0; i < 13; i++)
		if (v[i][0] == c)
			return i + 2;
	return 0;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	char c, lc, rc;
	unsigned l, r;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		l = rank(c);
		if (l == 10)
			fseek(fp, 1, SEEK_CUR);
		lc = getc(fp);
		fseek(fp, 1, SEEK_CUR);
		c = getc(fp);
		r = rank(c);
		if (r == 10)
			fseek(fp, 1, SEEK_CUR);
		rc = getc(fp);
		fseek(fp, 3, SEEK_CUR);
		c = getc(fp);
		if (lc == c)
			l += 16;
		if (rc == c)
			r += 16;
		if (l > r)
			printf("%s%c\n", v[(l & 15) - 2], lc);
		else if (l < r)
			printf("%s%c\n", v[(r & 15) - 2], rc);
		else
			printf("%s%c %s%c\n", v[(l & 15) - 2], lc,
			       v[(r & 15) - 2], rc);
		fseek(fp, 1, SEEK_CUR);
	}
	return 0;
}
