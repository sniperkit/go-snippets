#include <regex.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	regex_t email_regex;
	unsigned sbs = 48, n = 0;
	char c;
	char *sb = malloc(sbs);

	regcomp(&email_regex, "^((\"[0-9A-Za-z@.+-=]+\")|([0-9A-Za-z.+-=]+))"
			      "@\\w*(\\w+\\.)+\\w{2,4}$",
		REG_EXTENDED | REG_NOSUB);

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || n > 0) {
		if (c == '\n' || c == EOF) {
			sb[n] = '\0';
			puts(regexec(&email_regex, sb, 0, NULL, 0) ?
			     "false" : "true");
			n = 0;
		} else {
			if (n == sbs - 1) {
				sbs += sbs / 2;
				sb = realloc(sb, sbs);
			}
			sb[n++] = c;
		}
	}
	free(sb);
	return 0;
}
