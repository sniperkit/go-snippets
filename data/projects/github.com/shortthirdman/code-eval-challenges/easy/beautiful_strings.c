#include <stdio.h>
#include <stdlib.h>

static int cmpint(const void *p1, const void *p2) {
	return *(int *)p1 - *(int *)p2;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	char c;
	int ch[26] = { 0 }, i, s = -1;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || s >= 0) {
		if (s == -1)
			s = 0;
		if (c == '\n' || c == EOF) {
			qsort(&ch, 26, sizeof(int), cmpint);
			for (i = 25; ch[i] > 0 && i >= 0; i--) {
				s += (i + 1) * ch[i];
				ch[i] = 0;
			}
			printf("%d\n", s);
			s = -1;
		} else {
			c &= 223;
			if (c >= 'A' && c <= 'Z')
				ch[c - 'A']++;
		}
	}
	return 0;
}
