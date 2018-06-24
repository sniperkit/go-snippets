#include <stdbool.h>
#include <stdio.h>

#define READ  4
#define WRITE 2
#define GRANT 1

static char d(char s) {
	return s - '1';
}

int main(int argc, char *argv[]) {
	FILE *fp;
	char u, f, a;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "user_%c=>file_%c=>%c", &u, &f, &a) != EOF) {
		char uperm[18] = { 7, 3, 0, 6, 2, 4, 5, 1, 5,
				   3, 7, 1, 6, 0, 2, 4, 2, 6 };
		bool v = true;
		do {
			switch (a) {
			case 'r':
				fscanf(fp, "ead");
				if ((uperm[d(u) * 3 + d(f)] & READ) == 0)
					v = false;
				break;
			case 'w':
				fscanf(fp, "rite");
				if ((uperm[d(u) * 3 + d(f)] & WRITE) == 0)
					v = false;
				break;
			case 'g':
				fscanf(fp, "rant=>%c", &a);
				if ((uperm[d(u) * 3 + d(f)] & GRANT) == 0)
					v = false;
				switch (a) {
				case 'r':
					fscanf(fp, "ead=>user_%c", &u);
					uperm[d(u) * 3 + d(f)] |= READ;
					break;
				case 'w':
					fscanf(fp, "rite=>user_%c", &u);
					uperm[d(u) * 3 + d(f)] |= WRITE;
					break;
				case 'g':
					fscanf(fp, "rant=>user_%c", &u);
					uperm[d(u) * 3 + d(f)] |= GRANT;
				}
			}
		} while (getc(fp) == ' ' &&
			 fscanf(fp, "user_%c=>file_%c=>%c", &u, &f, &a) != EOF);
		puts(v ? "True" : "False");
	}
	return 0;
}
