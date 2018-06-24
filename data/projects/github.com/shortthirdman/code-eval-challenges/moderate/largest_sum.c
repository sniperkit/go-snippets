#include <limits.h>
#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int a, l = 0, m = INT_MIN;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d", &a) != EOF || m > INT_MIN) {
		char c = getc(fp);

		if (a + l > m)
			m = a + l;
		if (a + l > 0)
			l += a;
		else
			l = 0;
		if (c == '\n' || c == EOF) {
			printf("%d\n", m);
			l = 0;
			m = INT_MIN;
		}
	}
	return 0;
}
