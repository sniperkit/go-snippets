#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

struct node {
	struct node	*next;
	char		data;
	bool		uniq;
};

static void firstnrc(struct node *head) {
	struct node *curr = head;

	if (head) {
		while (curr) {
			if (curr->uniq) {
				printf("%c\n", curr->data);
				return;
			}
			curr = curr->next;
		}
	}
	printf("\n");
}

int main(int argc, char *argv[]) {
	FILE *fp;
	char c;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF) {
		struct node *head = NULL;
		while (c != '\n' && c != EOF) {
			struct node *temp;
			temp = malloc(sizeof(struct node));
			temp->data = c;
			temp->uniq = true;
			temp->next = NULL;
			if (head) {
				struct node *curr = head, *oldr = NULL;
				while (curr) {
					if (curr->data == temp->data) {
						curr->uniq = false;
						goto skip;
					}
					oldr = curr;
					curr = curr->next;
				}
				oldr->next = temp;
			} else {
				head = temp;
			}

skip:
			c = getc(fp);
		}
		firstnrc(head);
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
