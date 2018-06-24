#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int i, p = 0, s = 0, sbs = 32;
	char c;
	char *sb = calloc(sbs, 1);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || s > 0) {
		switch (c) {
		case '\n':
		case EOF:
			for (i = 0; i <= s; i++)
				sb[i] = 0;
			s = 0;
			p = 0;
			putchar('\n');
			break;
		case '+':
			sb[p]++;
			break;
		case '-':
			sb[p]--;
			break;
		case '.':
			printf("%c", sb[p]);
			break;
		case ',':
			sb[p] = getchar();
			break;
		case '<':
			if (p > 0) {
				p--;
			} else {
				if (++s == sbs) {
					sbs += sbs / 2;
					sb = realloc(sb, sbs);
				}
				for (i = s - 1; i > 0; i--)
					sb[i] = sb[i - 1];
				sb[0] = 0;
			}
			break;
		case '>':
			if (p == sbs) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
			}
			p++;
			if (p > s)
				s = p;
			break;
		case '[':
			if (sb[p] == 0) {
				i = 0;
				do {
					while ((c = getc(fp)) != ']' && c != '[');
					if (c == '[')
						i++;
					else
						i--;
				} while (i >= 0);
			}
			break;
		case ']':
			if (sb[p] != 0) {
				i = 0;
				do {
					fseek(fp, -2, SEEK_CUR);
					while ((c = getc(fp)) != '[' && c != ']') {
						fseek(fp, -2, SEEK_CUR);
					}
					if (c == ']')
						i++;
					else
						i--;
				} while (i >= 0);
			}
			break;
		}
	}
	free(sb);
	return 0;
}
