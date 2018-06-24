#include <stdio.h>

static const char *cat[] = { "This program is for humans",
			     "Still in Mama's arms",
			     "Preschool Maniac",
			     "Elementary school",
			     "Middle school",
			     "High school",
			     "College",
			     "Working for the man",
			     "The Golden Years" };
static const int age[] = { 0, 3, 5, 12, 15, 19, 23, 66, 101 };

int main(int argc, char *argv[]) {
	FILE *fp;
	int a;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d", &a) != EOF) {
		int c = 0;
		while (c < 9 && a >= age[c])
			c++;
		puts(cat[c % 9]);
	}
	return 0;
}
