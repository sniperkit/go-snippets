#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

struct item {
	unsigned num;
	unsigned pos;
};

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned a, i = 0, j = 1, ibs = 32, k, m;
	struct item *ib = malloc(ibs * sizeof(struct item));

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%u", &a) != EOF) {
		bool f = false;
		m = a;
		ib[i].num = a;
		ib[i++].pos = j++;
		while (getc(fp) == ' ') {
			fscanf(fp, "%u", &a);
			for (k = 0; k < i; k++) {
				if (ib[k].num == a) {
					ib[k].pos = 0;
					break;
				}
			}
			if (k == i) {
				if (m < a)
					m = a;
				if (i == ibs) {
					ibs += ibs / 2;
					ib = realloc(ib, ibs * sizeof(struct item));
				}
				ib[i].num = a;
				ib[i++].pos = j;
			}
			j++;
		}
		for (k = 0; k < i; k++) {
			if (ib[k].pos != 0 && ib[k].num <= m) {
				m = ib[k].num;
				j = ib[k].pos;
				f = true;
			}
		}
		if (!f)
			j = 0;
		printf("%u\n", j);
		i = 0;
		j = 1;
	}
	free(ib);
	return 0;
}
