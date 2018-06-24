#include <stdbool.h>
#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	bool p = true;
	char c;
	int count = 0, arrow = 0;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || !p) {
		if (c == '\n' || c == EOF) {
			printf("%d\n", count);
			count = 0;
			arrow = 0;
			p = true;
			continue;
		}
		p = false;
		switch (arrow) {
		case 0:
			switch (c) {
			case '>':
				arrow++;
				break;
			case '<':
				arrow--;
				break;
			}
			break;
		case 1:
			switch (c) {
			case '>':
				arrow++;
				break;
			case '<':
				arrow = -1;
				break;
			default:
				arrow = 0;
			}
			break;
		case -1:
			switch (c) {
			case '-':
				arrow--;
				break;
			case '>':
				arrow = 1;
				break;
			}
			break;
		case 2:
			switch (c) {
			case '-':
				arrow++;
				break;
			case '<':
				arrow = -1;
				break;
			}
			break;
		case -2:
			switch (c) {
			case '-':
				arrow--;
				break;
			case '>':
				arrow = 1;
				break;
			default:
				arrow = -1;
			}
			break;
		case 3:
			switch (c) {
			case '-':
				arrow++;
				break;
			case '<':
				arrow = -1;
				break;
			default:
				arrow = 1;
			}
			break;
		case -3:
			switch (c) {
			case '<':
				arrow--;
				break;
			case '>':
				arrow = 1;
				break;
			default:
				arrow = 0;
			}
			break;
		case 4:
			switch (c) {
			case '>':
				count++;
				arrow = 1;
				break;
			case '<':
				arrow = -1;
				break;
			default:
				arrow = 0;
			}
			break;
		case -4:
			switch (c) {
			case '<':
				count++;
				arrow = -1;
				break;
			case '>':
				arrow = 1;
				break;
			default:
				arrow = -2;
			}
			break;
		}
	}
	return 0;
}
