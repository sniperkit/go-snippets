#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

struct node {
	struct node	*next;
	int		data;
};

int main(int argc, char *argv[]) {
	FILE *fp;
	int n;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d;", &n) != EOF) {
		struct node *head = NULL;
		int a;
		bool done = false;
		while (n > 0) {
			fscanf(fp, "%d", &a);
			if (done)
				goto skip;

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
					done = true;
					goto skip;
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

skip:
			fseek(fp, 1, SEEK_CUR);
			n--;
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
