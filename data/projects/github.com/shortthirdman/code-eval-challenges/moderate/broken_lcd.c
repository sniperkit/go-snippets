#include <stdbool.h>
#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	char c, dig[13];
	int lcd[12], d[12], n = 0;
	const int digits[] = { 252,  96, 218, 242, 102,
			       182, 190, 224, 254, 246 };

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		int i, j;
		bool f = false;
		for (i = 0; i < 12; i++) {
			lcd[i] = 0;
			for (j = 7; j >= 0; j--) {
				lcd[i] += (c == '1') << j;
				c = getc(fp);
			}
			c = getc(fp);
		}
		while (c != '\n' && c != EOF) {
			dig[n++] = c;
			c = getc(fp);
		}
		d[0] = digits[dig[0] - '0'];
		j = 1;
		for (i = 1; i < n; i++)
			if (dig[i] == '.')
				d[j - 1]++;
			else
				d[j++] = digits[dig[i] - '0'];
		for (i = 0; i < 14 - n && !f; i++) {
			f = true;
			for (j = 0; j < n - 1; j++) {
				if ((lcd[i + j] & d[j]) != d[j]) {
					f = false;
					break;
				}
			}
		}
		printf("%d\n", f);
		n = 0;
	}
	return 0;
}
