#include <stdio.h>
#include <stdlib.h>

struct stack {
	unsigned length;
	unsigned size;
	int      *ib;
};

void push(struct stack *, int);
int pop(struct stack *);

void push(struct stack *st, int i) {
	if (st->length == st->size) {
		st->size += st->size / 2;
		st->ib = realloc(st->ib, st->size * sizeof(int));
	}
	st->ib[st->length++] = i;
	return;
}

int pop(struct stack *st) {
	if (st->length == 0)
		return -1;
	return st->ib[--st->length];
}

int main(int argc, char *argv[]) {
	FILE *fp;
	int a;
	char c;
	struct stack *st = malloc(sizeof(struct stack));
	st->size = 46;
	st->ib = malloc(st->size * sizeof(int));

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d", &a) != EOF) {
		st->length = 0;
		push(st, a);
		while ((c = getc(fp)) == ' ') {
			fscanf(fp, "%d", &a);
			push(st, a);
		}
		while (st->length) {
			printf("%d", pop(st));
			pop(st);
			if (st->length)
				putchar(' ');
		}
		putchar('\n');
	}
	free(st->ib);
	free(st);
	return 0;
}
