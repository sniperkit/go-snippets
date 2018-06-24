#include <stdio.h>
#include <stdlib.h>

static int cmpint(const void *p1, const void *p2) {
	return *(int *)p2 - *(int *)p1;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned time[20], n, h, m, s, i;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d:%d:%d", &h, &m, &s) != EOF) {
		n = 1;
		time[0] = (h << 12) + (m << 6) + s;
		while (getc(fp) == ' ') {
			fscanf(fp, "%d:%d:%d", &h, &m, &s);
			time[n++] = (h << 12) + (m << 6) + s;
		}
		qsort(&time, n, sizeof(int), cmpint);
		printf("%02d:%02d:%02d", time[0] >> 12, (time[0] >> 6) & 63, time[0] & 63);
		for (i = 1; i < n; i++)
			printf(" %02d:%02d:%02d", time[i] >> 12, (time[i] >> 6) & 63, time[i] & 63);
		putchar('\n');
	}
	return 0;
}
