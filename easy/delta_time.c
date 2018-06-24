#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int h1, m1, s1, h2, m2, s2;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d:%d:%d %d:%d:%d\n", &h1, &m1, &s1, &h2, &m2, &s2) != EOF) {
		int t = abs((h1 - h2) * 3600 + (m1 - m2) * 60 + s1 - s2);
		printf("%02d:%02d:%02d\n", t / 3600, (t / 60) % 60, t % 60);
	}
	return 0;
}
