#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

static int next_permu(char *c) {
	char t;
	int i, k = -1, l = 0;

	for (i = 3; i >= 0; i--) {
		if (c[i] < c[i + 1]) {
			k = i;
			break;
		}
	}
	if (k == -1)
		return 0;
	for (i = 4; i > 0; i--) {
		if (c[i] > c[k]) {
			l = i;
			break;
		}
	}

	t = c[k];
	c[k] = c[l];
	c[l] = t;

	l = 4 - k;
	for (i = 0; (l - 1) - i > i; i++) {
		t = c[4 - i];
		c[4 - i] = c[k + 1 + i];
		c[k + 1 + i] = t;
	}
	return l;
}

static int ccmp(const void *b1, const void *b2) {
	char *c1 = (char *)b1, *c2 = (char *)b2;
	if (*c1 < *c2)
		return -1;
	return *c1 > *c2;
}

static bool gives42(char *s) {
	int o1, o2, o3, o4, r1, r2, r3, r4;
	for (o1 = 0; o1 < 3; o1++) {
		r1 = (o1 == 0) ? s[0] + s[1] : (o1 == 1) ? s[0] - s[1] : s[0] * s[1];
		for (o2 = 0; o2 < 3; o2++) {
			r2 = (o2 == 0) ? r1 + s[2] : (o2 == 1) ? r1 - s[2] : r1 * s[2];
			for (o3 = 0; o3 < 3; o3++) {
				r3 = (o3 == 0) ? r2 + s[3] : (o3 == 1) ? r2 - s[3] : r2 * s[3];
				for (o4 = 0; o4 < 3; o4++) {
					r4 = (o4 == 0) ? r3 + s[4] : (o4 == 1) ? r3 - s[4] : r3 * s[4];
					if (r4 == 42)
						return true;
				}
			}
		}
	}
	return false;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	char s[5];

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%hhd %hhd %hhd %hhd %hhd", &s[0], &s[1], &s[2], &s[3], &s[4]) != EOF) {
		bool f = false;
		qsort(s, 5, 1, ccmp);
		f = gives42(s);
		while (!f && next_permu(s))
			f = gives42(s);
		puts(f ? "YES" : "NO");
	}
	return 0;
}
