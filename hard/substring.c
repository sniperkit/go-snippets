#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

static bool subcheck(char *sb, int n1, int n2, int i, int j) {
	while (i < n1 && j < n2) {
		if (sb[n1 + j] == '*') {
			bool a = false;
			int k;
			for (k = i; k < n1; k++) {
				if (subcheck(sb, n1, n2, k, j + 1)) {
					a = true;
					break;
				}
			}
			return a;
		}
		if (sb[i] != sb[n1 + j])
			return false;
		i++;
		j++;
	}
	if (j < n2)
		return false;
	return true;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	bool a = true, b = false;
	int sbs = 32, s = 0, n1 = -1, n2 = -1, i;
	char c;
	char *sb = malloc(sbs);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || s > 0) {
		switch (c) {
		case ',':
			n1 = s;
			break;
		case '\n':
		case EOF:
			while (s >= n1 && sb[s - 1] == '*')
				s--;
			n2 = s - n1;
			b = false;
			for (i = 0; i < n1; i++) {
				if (subcheck(sb, n1, n2, i, 0)) {
					b = true;
					break;
				}
			}
			puts(b ? "true" : "false");
			s = 0;
			n1 = -1;
			n2 = -1;
			a = true;
			b = false;
			break;
		default:
			if (n1 >= 0) {
				if (a) {
					if (c == '*')
						break;
					a = false;
				}
				if (c == '\\') {
					b = !b;
					if (b)
						break;
				} else if (c == '*') {
					if (b) {
						c = '+';
						b = false;
					} else {
						a = true;
					}
				}
			} else if (c == '*') {
				c = '+';
			}
			if (s == sbs - 1) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
			}
			sb[s++] = c;
		}
	}
	free(sb);
	return 0;
}
