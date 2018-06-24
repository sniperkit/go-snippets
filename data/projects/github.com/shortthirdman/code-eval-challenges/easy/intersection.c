#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

struct node {
	int		data;
	struct node	*next;
};

int main(int argc, char *argv[]) {
	FILE *fp;
	int a;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d", &a) != EOF) {
		struct node *head = NULL, *temp, *tail;
		bool printing = false, done = false;
		temp = malloc(sizeof(struct node));
		temp->data = a;
		temp->next = head;
		head = temp;
		tail = head;
		while (getc(fp) == ',') {
			fscanf(fp, "%d", &a);
			temp = malloc(sizeof(struct node));
			temp->data = a;
			temp->next = tail->next;
			tail->next = temp;
			tail = temp;
		}
		fscanf(fp, "%d", &a);
		temp = head;
		while (temp->data < a) {
			tail = temp;
			if (!temp->next)
				break;
			temp = temp->next;
		}
		if (temp->data == a) {
			printing = true;
			printf("%d", a);
		} else if (temp->data < a) {
			done = true;
		}
		while (getc(fp) == ',') {
			fscanf(fp, "%d", &a);
			if (!done) {
				while (temp->data < a) {
					tail = temp;
					if (!temp->next)
						break;
					temp = temp->next;
				}
				if (temp->data == a) {
					if (printing)
						printf(",");
					else
						printing = true;
					printf("%d", a);
				} else if (temp->data < a) {
					done = true;
				}
			}
		}
		printf("\n");
		while (head) {
			temp = head;
			head = head->next;
			free(temp);
		}
	}
	return 0;
}
