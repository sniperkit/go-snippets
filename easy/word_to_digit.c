#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	char c;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%c", &c) != EOF) {
		switch (c) {
		case 'z':
			printf("0");
			fseek(fp, 3, SEEK_CUR);
			break;
		case 'o':
			printf("1");
			fseek(fp, 2, SEEK_CUR);
			break;
		case 't':
			fscanf(fp, "%c", &c);
			if (c == 'w') {
				printf("2");
				fseek(fp, 1, SEEK_CUR);
			} else {
				printf("3");
				fseek(fp, 3, SEEK_CUR);
			}
			break;
		case 'f':
			fscanf(fp, "%c", &c);
			if (c == 'o')
				printf("4");
			else
				printf("5");
			fseek(fp, 2, SEEK_CUR);
			break;
		case 's':
			fscanf(fp, "%c", &c);
			if (c == 'i') {
				printf("6");
				fseek(fp, 1, SEEK_CUR);
			} else {
				printf("7");
				fseek(fp, 3, SEEK_CUR);
			}
			break;
		case 'e':
			printf("8");
			fseek(fp, 4, SEEK_CUR);
			break;
		case 'n':
			printf("9");
			fseek(fp, 3, SEEK_CUR);
			break;
		}
		if (fscanf(fp, "%c", &c) == EOF || c == '\n')
			printf("\n");
	}
	return 0;
}
