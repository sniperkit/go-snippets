#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>

struct link {
	short from;
	short to;
};

static int to(struct link *p, int n, int a) {
	int i;
	for (i = 0; i < n; i++) {
		if (p[i].from == p[a].to) {
			return i;
		}
	}
	return -1;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	struct link chain[500];
	char c, s[6];

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = fgetc(fp)) != EOF) {
		bool good = true;
		int n = 0, i = 0, b = -1;
		while (c != EOF && c != '\n') {
			if (n >= 500)
				good = false;
			s[i++] = c;
			while ((c = fgetc(fp)) != '-') {
				if (i > 4) {
					i = 0;
					good = false;
				}
				s[i++] = c;
			}
			s[i] = '\0';
			if (good) {
				if (s[0] == 'B') {
					chain[n].from = 0;
					b = n;
				} else if (s[0] == 'E') {
					good = false;
				} else {
					sscanf(s, "%hd", &chain[n].from);
				}
			}
			i = 0;
			while ((c = fgetc(fp)) != ';' && c != '\n' && c != EOF) {
				if (i > 4) {
					i = 0;
					good = false;
				}
				s[i++] = c;
			}
			s[i] = '\0';
			if (good) {
				if (s[0] == 'E')
					chain[n].to = 1;
				else if (s[0] == 'B')
					good = false;
				else
					sscanf(s, "%hd", &chain[n].to);
				if (chain[n].from == chain[n].to)
					good = false;
			}
			i = 0;
			n++;
			if (c == ';')
				c = fgetc(fp);
		}
		if (b == -1)
			good = false;
		for (i = 0; good && i < n - 1; i++) {
			b = to(chain, n, b);
			if (b == -1 || (chain[b].to == 1 && i < n - 2))
				good = false;
		}
		puts(good && chain[b].to == 1 ? "GOOD" : "BAD");
		n = 0;
		b = -1;
	}
	return 0;
}
