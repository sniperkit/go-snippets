#include <stdio.h>
#include <stdlib.h>

static int cmpint(const void *p1, const void *p2) {
	return *(int *)p1 - *(int *)p2;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	int s1, s2, s3, s4, s5, s6;
	unsigned h[2], b[3], r[15], i, j;
	char c;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "[%d,%d] [%d,%d]|", &s1, &s2, &s3, &s4) != EOF) {
		j = 0;
		h[0] = (unsigned)abs(s1 - s3);
		h[1] = (unsigned)abs(s2 - s4);
		qsort(&h, 2, sizeof(unsigned), cmpint);
		while ((c = getc(fp)) == '(') {
			fscanf(fp, "%u [%d,%d,%d] [%d,%d,%d]);",
			       &i, &s1, &s2, &s3, &s4, &s5, &s6);
			b[0] = (unsigned)abs(s1 - s4);
			b[1] = (unsigned)abs(s2 - s5);
			b[2] = (unsigned)abs(s3 - s6);
			qsort(&b, 3, sizeof(int), cmpint);
			if (b[0] <= h[0] && b[1] <= h[1]) { /* || diagonal(b, h) */
				r[j++] = i;
			}
		}
		if (j == 0) {
			puts("-");
		} else {
			qsort(&r, j, sizeof(unsigned), cmpint);
			printf("%u", r[0]);
			for (i = 1; i < j; i++)
				printf(",%d", r[i]);
			putchar('\n');
		}
	}
	return 0;
}
