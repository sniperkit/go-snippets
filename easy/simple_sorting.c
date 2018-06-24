#include <stdio.h>
#include <stdlib.h>

struct node {
	float		num;
	struct node	*next;
};

int main(int argc, char *argv[]) {
	FILE *fp;
	float a;
	struct node *head = NULL, *temp, *curr;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%f", &a) != EOF) {
		temp = malloc(sizeof(struct node));
		temp->num = a;
		temp->next = head;
		head = temp;
		while (getc(fp) == ' ') {
			if (fscanf(fp, "%f", &a) == EOF)
				goto fail_eof;
			temp = malloc(sizeof(struct node));
			temp->num = a;
			if (a < head->num) {
				temp->next = head;
				head = temp;
			} else {
				curr = head;
				while (curr->next && curr->next->num < a)
					curr = curr->next;
				temp->next = curr->next;
				curr->next = temp;
			}
		}
		while (head) {
			printf("%.3f", head->num);
			if (head->next) {
				printf(" ");
			}
			curr = head;
			head = head->next;
			free(curr);
		}
		printf("\n");
	}
	return 0;
fail_eof:
	printf("unexpected end of input\n");
	while (head) {
		curr = head;
		head = head->next;
		free(curr);
	}
	return EXIT_FAILURE;
}
