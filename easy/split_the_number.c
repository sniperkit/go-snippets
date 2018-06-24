#include <stdbool.h>
#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	char c, z[10];

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		int n = 0, i;
		while (c >= '0' && c <= '9') {
			z[n] = c - '0';
			n++;
			c = getc(fp);
		}
		c = getc(fp);
		int num = 0, op = 0;
		bool sg = true;
		for (i = 0; i < n;) {
			if (c >= 'a' && c <= 'z') {
				op = 10 * op + z[i];
				i++;
			} else {
				if (sg)
					num += op;
				else
					num -= op;
				op = 0;
				sg = c == '+';
			}
			c = getc(fp);
		}
		if (sg)
			num += op;
		else
			num -= op;
		printf("%d\n", num);
	}
	return 0;
}
