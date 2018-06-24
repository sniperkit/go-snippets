#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

struct node {
	int		x;
	int		y;
	struct node	*next;
};

static int zsum(int a, int b) {
	int z = 0;
	while (a) {
		z += a % 10;
		a /= 10;
	}
	while (b) {
		z += b % 10;
		b /= 10;
	}
	return z;
}

static bool unlisted(int x, int y, struct node *h1, struct node *h2) {
	struct node *temp = h1;
	if (temp) {
		while (temp->x != x || temp->y != y) {
			if (!temp->next)
				break;
			temp = temp->next;
		}
		if (temp->x == x && temp->y == y)
			return false;
	}
	temp = h2;
	if (temp) {
		while (temp->x != x || temp->y != y) {
			if (!temp->next)
				break;
			temp = temp->next;
		}
		if (temp->x == x && temp->y == y)
			return false;
	}
	return true;
}

int main(void) {
	int k = 0, s = 19;
	struct node *h1 = malloc(sizeof(struct node));
	struct node *h2 = NULL;
	h1->x = 0;
	h1->y = 0;
	h1->next = NULL;

	while (h1) {
		struct node *temp = h1;
		h1 = h1->next;
		temp->next = h2;
		h2 = temp;

		if (zsum((h2->x) + 1, h2->y) <= s &&
		    unlisted((h2->x) + 1, h2->y, h1, h2)) {
			temp = malloc(sizeof(struct node));
			temp->x = (h2->x) + 1;
			temp->y = h2->y;
			temp->next = h1;
			h1 = temp;
		}
		if (zsum(h2->x, (h2->y) + 1) <= s &&
		    unlisted(h2->x, (h2->y) + 1, h1, h2)) {
			temp = malloc(sizeof(struct node));
			temp->x = h2->x;
			temp->y = (h2->y) + 1;
			temp->next = h1;
			h1 = temp;
		}
	}

	while (h2) {
		struct node *temp = h2;
		if (h2->y > 0)
			k++;
		h2 = h2->next;
		free(temp);
	}

	printf("%d\n", k * 4 + 1);
	return 0;
}
