#include <stdio.h>

static int d0(char *s) {
	switch (s[0]) {
	case 'o':
		return 1;
	case 't':
		if (s[1] == 'w' && s[2] == 'o')
			return 2;
		else if (s[2] == 'r')
			return 3;
		break;
	case 'f':
		if (s[4] != '\0')
			break;
		else if (s[1] == 'o')
			return 4;
		else
			return 5;
	case 's':
		if (s[3] == '\0')
			return 6;
		else if (s[1] == 'e' && s[5] == '\0')
			return 7;
		break;
	case 'e':
		if (s[5] == '\0')
			return 8;
		break;
	case 'n':
		if (s[4] == '\0')
			return 9;
	}
	return -1;
}

static int d1(char *s) {
	switch (s[0]) {
	case 't':
		if (s[3] == '\0')
			return 10;
		else if (s[6] == 'e')
			return 13;
		else if (s[5] == 'e')
			return 12;
		break;
	case 'e':
		if (s[1] == 'l')
			return 11;
		else if (s[5] == 'e')
			return 18;
		break;
	case 'f':
		if (s[4] == 'e')
			return 15;
		else if (s[4] == 't')
			return 14;
		break;
	case 's':
		if (s[1] == s[6])
			return 17;
		else if (s[4] == s[5])
			return 16;
		break;
	case 'n':
		if (s[5] == 'e')
			return 19;
	}
	return -1;
}

static int d2(char *s) {
	switch (s[0]) {
	case 't':
		if (s[2] == 'i')
			return 30;
		else if (s[1] == 'w')
			return 20;
		break;
	case 'f':
		if (s[1] == 'i')
			return 50;
		return 40;
	case 's':
		if (s[1] == 'e')
			return 70;
		return 60;
	case 'e':
		return 80;
	case 'n':
		return 90;
	}
	return -1;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	char s[10], c;
	int i = 0, n = 1, h = 0, q = 0, r = 0, d;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || i > 0) {
		if (c == ' ' || c == '\n' || c == EOF) {
			s[i] = '\0';
			if (s[0] == 'n' && s[1] == 'e') {
				n = -1;
			} else if (s[0] == 'z') {
				q = 0;
			} else {
				d = d0(s);
				if (d >= 0) {
					h = d;
				} else {
					d = d1(s);
					if (d >= 0) {
						q += d;
					} else {
						d = d2(s);
						if (d >= 0) {
							q += d;
						} else if (s[0] == 'h') {
							q = h * 100;
							h = 0;
						} else if (s[0] == 't') {
							r += (h + q) * 1000;
							h = 0;
							q = 0;
						} else if (s[0] == 'm') {
							r = (h + q) * 1000000;
							h = 0;
							q = 0;
						}
					}
				}
			}
			if (c != ' ') {
				printf("%d\n", n * (r + q + h));
				n = 1;
				h = 0;
				q = 0;
				r = 0;
			}
			i = 0;
		} else {
			s[i++] = c;
		}
	}
	return 0;
}
