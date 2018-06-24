#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	int x1, y1, x2, y2;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d %d %d %d", &x1, &y1, &x2, &y2) != EOF) {
		puts((x1 == x2) ? ((y1 == y2) ? "here" :
				   (y1 < y2) ? "N" : "S") :
		     (y1 == y2) ? ((x1 < x2) ? "E" : "W") :
		     ((x1 < x2) ? ((y1 < y2) ? "NE" : "SE") :
				  ((y1 < y2) ? "NW" : "SW")));
	}
	return 0;
}
