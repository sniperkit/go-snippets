#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

static void trim_right(char *p) {
	int i = 0;
	while (p[i] != '\0')
		i++;
	while (p[i] > 0 && p[i - 1] != '(' && p[i - 1] != ')')
		i--;
	p[i] = '\0';
}

static bool is_balanced(char *p, unsigned c) {
	if (*p == '\0')
		return !c;
	if (p[0] != ':' && p[0] != '(' && p[0] != ')')
		return is_balanced(p + 1, c);
	trim_right(p);
	switch (p[0]) {
	case '(':
		return is_balanced(p + 1, c + 1);
	case ')':
		if (c)
			return is_balanced(p + 1, c - 1);
		return false;
	case ':':
		switch (p[1]) {
		case '(':
			return is_balanced(p + 2, c) ||
			       is_balanced(p + 2, c + 1);
		case ')':
			if (c)
				return is_balanced(p + 2, c) ||
				       is_balanced(p + 2, c - 1);
			return is_balanced(p + 2, c);
		default:
			return is_balanced(p + 1, c);
		}
	default:
		return false;
	}
}

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned s = 0, sbs = 32;
	char c;
	char *sb = malloc(sbs);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || s > 0) {
		if (c == '\n' || c == EOF) {
			sb[s] = '\0';
			puts(is_balanced(sb, 0) ? "YES" : "NO");
			s = 0;
		} else {
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
