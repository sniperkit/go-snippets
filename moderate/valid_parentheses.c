#include <stdio.h>
#include <stdlib.h>

struct node {
	char		sym;
	struct node	*next;
};

int main(int argc, char *argv[]) {
	FILE *fp;
	char c;
	struct node *head = NULL, *temp = NULL;
	char valid = -1;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((c = getc(fp)) != EOF || valid > -1) {
		if (valid == 0) {
			if (c == '\n' || c == EOF) {
				valid = -1;
			}
			continue;
		}
		if (valid == -1)
			valid = 1;
		switch (c) {
		case '\n':
		case EOF:
			if (head) {
				puts("False");
				while (head) {
					struct node *done = head;
					head = head->next;
					free(done);
				}
				valid = -1;
			} else {
				puts("True");
			}
			break;
		case ')':
			if (head && head->sym == '(') {
				struct node *done = head;
				head = head->next;
				free(done);
			} else {
				valid = 0;
			}
			break;
		case ']':
			if (head && head->sym == '[') {
				struct node *done = head;
				head = head->next;
				free(done);
			} else {
				valid = 0;
			}
			break;
		case '}':
			if (head && head->sym == '{') {
				struct node *done = head;
				head = head->next;
				free(done);
			} else {
				valid = 0;
			}
			break;
		default:
			temp = malloc(sizeof(struct node));
			temp->sym = c;
			temp->next = head;
			head = temp;
		}
		if (valid == 0) {
			puts("False");
			while (head) {
				struct node *done = head;
				head = head->next;
				free(done);
			}
		}
	}
	return 0;
}
