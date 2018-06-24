#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

struct team {
	unsigned id;
	unsigned countries;
};

static int cmpt(const void *p1, const void *p2) {
	struct team *t1 = (struct team *)p1;
	struct team *t2 = (struct team *)p2;
	return (int)(t1->id - t2->id);
}

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned a, i, j, c = 0, t = 0, tbs = 16;
	struct team *tb = malloc(tbs * sizeof(struct team));
	char ch;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%u", &a) != EOF) {
		while (true) {
			while (true) {
				for (i = 0; i < t; i++)
					if (tb[i].id == a)
						break;
				if (i == t) {
					if (t == tbs) {
						tbs += tbs / 2;
						tb = realloc(tb, tbs * sizeof(struct team));
					}
					tb[t].id = a;
					tb[t].countries = 1 << c;
					t++;
				} else {
					tb[i].countries |= 1 << c;
				}
				while ((ch = getc(fp)) == ' ') ;
				if (ch == '|') {
					break;
				}
				if (ch < '1' || ch > '9')
					break;
				fseek(fp, -1, SEEK_CUR);
				fscanf(fp, "%u", &a);
			}
			if (ch != '|') {
				break;
			}
			c++;
			fseek(fp, 1, SEEK_CUR);
			fscanf(fp, "%u", &a);
		}

		qsort(tb, t, sizeof(struct team), cmpt);
		for (i = 0; i < t; i++) {
			if (i)
				putchar(' ');
			j = 0;
			while (tb[i].countries && !(tb[i].countries & (1 << j)))
				j++;
			tb[i].countries &= ~(1 << j);
			printf("%u:%u", tb[i].id, j + 1);
			while (tb[i].countries) {
				while (!(tb[i].countries & (1 << j)))
					j++;
				tb[i].countries &= ~(1 << j);
				printf(",%u", j + 1);
			}
			putchar(';');
		}
		putchar('\n');
		c = 0;
		t = 0;
	}
	free(tb);
	return 0;
}
