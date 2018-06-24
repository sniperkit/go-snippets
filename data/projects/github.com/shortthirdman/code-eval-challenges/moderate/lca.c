#include <stdio.h>
#include <stdlib.h>

struct node {
	struct node	*next;
	int		data;
};

int main(int argc, char *argv[]) {
	FILE *fp;
	int a, b;
	/*
	 *	node -> parent
	 *	 3	 8
	 *	 8	30
	 *	10	20
	 *	20	 8
	 *	29	20
	 *	30	 0
	 *	52	30
	 */
	int t[53] = { 0, 0, 0, 8, 0, 0, 0, 0,30, 0,
		     20, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		      8, 0, 0, 0, 0, 0, 0, 0, 0,20,
		      0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		      0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		      0, 0,30 };

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d %d", &a, &b) != EOF) {
		struct node *head = NULL;
		while (1) {
			if (a) {
				struct node *temp, *curr = head, *oldr = NULL;
				if (head) {
					while (curr->data < a) {
						oldr = curr;
						if (!curr->next)
							break;
						curr = curr->next;
					}
					if (curr->data == a) {
						printf("%d\n", a);
						break;
					}
				}
				temp = malloc(sizeof(struct node));
				temp->data = a;
				if (oldr) {
					temp->next = oldr->next;
					oldr->next = temp;
				} else {
					temp->next = head;
					head = temp;
				}
			}

			if (b) {
				struct node *temp, *curr = head, *oldr = NULL;
				while (curr->data < b) {
					oldr = curr;
					if (!curr->next)
						break;
					curr = curr->next;
				}
				if (curr->data == b) {
					printf("%d\n", b);
					break;
				}
				temp = malloc(sizeof(struct node));
				temp->data = b;
				if (oldr) {
					temp->next = oldr->next;
					oldr->next = temp;
				} else {
					temp->next = head;
					head = temp;
				}
			}

			a = t[a];
			b = t[b];
		}

		if (head) {
			struct node *curr = head, *oldr;
			while (curr) {
				oldr = curr;
				curr = curr->next;
				free(oldr);
			}
		}
	}
	return 0;
}
