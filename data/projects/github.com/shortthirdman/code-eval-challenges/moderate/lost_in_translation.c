#include <stdbool.h>
#include <stdio.h>

static const char *codel = "nug"
	"rbc vjnmkf kd yxyqci na rbc zjkfoscdd ew rbc ujllmcp"
	"tc rbkso rbyr ejp mysljylc kd kxveddknmc re jsicpdrysi"
	"de kr kd eoya kw aej icfkici re zjkr";
static const char *decol = "bjv"
	"the public is amazed by the quickness of the juggler"
	"we think that our language is impossible to understand"
	"so it is okay if you decided to quit";

int main(int argc, char *argv[]) {
	FILE *fp;
	bool p = true;
	char c;
	char m[26] = { 0 }, n[26] = { 0 };
	unsigned i = 0, j;

	while (codel[i] != '\0') {
		if (codel[i] >= 'a' && codel[i] <= 'z') {
			m[codel[i] - 'a'] = decol[i];
			n[decol[i] - 'a'] = codel[i];
		}
		i++;
	}
	for (i = 0; i < 26; i++) {
		if (m[i] == 0) {
			for (j = 0; j < 26; j++) {
				if (n[j] == 0) {
					m[i] = (char)j + 'a';
					n[j] = (char)i + 'a';
					break;
				}
			}
		}
	}

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		if (c >= 'a' && c <= 'z')
			putchar(m[c - 'a']);
		else
			putchar(c);
		p = c == '\n';
	}
	if (!p)
		putchar('\n');
	return 0;
}
