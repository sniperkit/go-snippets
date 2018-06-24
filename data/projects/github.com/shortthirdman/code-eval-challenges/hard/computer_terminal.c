#include <stdbool.h>
#include <stdio.h>

static void cls(char c[10][10]) {
	int i, j;
	for (i = 0; i < 10; i++)
		for (j = 0; j < 10; j++)
			c[i][j] = ' ';
}

static int prn(char c[10][10], int x, int y, char d, bool insr) {
	if (insr) {
		int i;
		for (i = 8; i >= y; i--)
			c[x][i + 1] = c[x][i];
	}
	c[x][y] = d;
	if (y < 9)
		y++;
	return y;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	bool ctrl = false, insr = false, move = false;
	char curr, scrn[10][10];
	int x = 0, y = 0, i, j;
	cls(scrn);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%c", &curr) != EOF) {
		if (curr == '\n')
			continue;
		if (ctrl) {
			switch (curr) {
			case '^':
				y = prn(scrn, x, y, curr, insr);
				ctrl = false;
				break;
			case 'c':
				cls(scrn);
				ctrl = false;
				break;
			case 'e':
				for (i = 8; i >= y; i--)
					scrn[x][i] = ' ';
				ctrl = false;
				break;
			case 'i':
				insr = true;
				ctrl = false;
				break;
			case 'o':
				insr = false;
				ctrl = false;
				break;
			case 'h':
				x = 0;
			case 'b':
				y = 0;
				ctrl = false;
				break;
			case 'd':
				if (x < 9)
					x++;
				ctrl = false;
				break;
			case 'u':
				if (x > 0)
					x--;
				ctrl = false;
				break;
			case 'l':
				if (y > 0)
					y--;
				ctrl = false;
				break;
			case 'r':
				if (y < 9)
					y++;
				ctrl = false;
				break;
			default:
				if (curr < '0' || curr > '9')
					break;
				if (!move) {
					move = true;
					x = curr - '0';
				} else {
					move = false;
					y = curr - '0';
					ctrl = false;
				}
			}
		} else {
			switch (curr) {
			case '^':
				ctrl = true;
				break;
			default:
				y = prn(scrn, x, y, curr, insr);
			}
		}
	}

	for (i = 0; i < 10; i++) {
		for (j = 0; j < 10; j++)
			putchar(scrn[i][j]);
		putchar('\n');
	}
	return 0;
}
