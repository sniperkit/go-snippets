#include <stdbool.h>
#include <stdio.h>

static const char *phrases[] = { ", yeah!", ", this is crazy, I tell ya.",
				 ", can U believe this?", ", eh?", ", aw yea.",
				 ", yo.", "? No way!", ". Awesome!" };

int main(int argc, char *argv[]) {
	FILE *fp;
	char c;
	bool p = true, l = false;
	int i = 0;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		if (c == '.' || c == '!' || c == '?') {
			if (l) {
				printf("%s", phrases[i]);
				l = false;
				i = (i + 1) % 8;
			} else {
				putchar(c);
				l = true;
			}
		} else {
			putchar(c);
		}
		p = c == '\n';
	}
	if (!p)
		putchar('\n');
	return 0;
}
