#include <stdio.h>

static const char *digits = "-**----*--***--***---*---****--**--****--**---**--"
			    "*--*--**-----*----*-*--*-*----*-------*-*--*-*--*-"
			    "*--*---*---**---**--****-***--***----*---**---***-"
			    "*--*---*--*-------*----*----*-*--*--*---*--*----*-"
			    "-**---***-****-***-----*-***---**---*----**---**--"
			    "--------------------------------------------------";
static const int w = 5, h = 6;

int main(int argc, char *argv[]) {
	FILE *fp;
	int i, j, k;
	char line[18];

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fgets(line, 18, fp) != 0)
		for (i = 0; i < h; i++) {
			for (j = 0; line[j] != '\n' && line[j] != '\0'; j++)
				for (k = 0; k < w; k++)
					if (line[j] >= '0' && line[j] <= '9')
						putchar(digits[i * 10 * w + (line[j] - '0') * w + k]);
			putchar('\n');
		}
	return 0;
}
