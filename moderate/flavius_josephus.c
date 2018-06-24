#include <stdio.h>
#include <stdlib.h>

struct node {
	int		num;
	struct node	*next;
};

int main(int argc, char *argv[]) {
	FILE *fp;
	int n, m;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d,%d", &n, &m) != EOF) {
		int i, j;
		struct node *head, *temp, *curr;
		temp = malloc(sizeof(struct node));
		temp->num = n - 1;
		temp->next = NULL;
		head = temp;
		curr = temp;
		for (i = n - 2; i >= 0; i--) {
			temp = malloc(sizeof(struct node));
			temp->num = i;
			temp->next = head;
			head = temp;
		}
		curr->next = head;
		for (i = 0; i < n; i++) {
			if (i)
				printf(" ");
			for (j = 0; j < m - 1; j++) {
				curr = curr->next;
			}
			printf("%d", curr->next->num);
			temp = curr->next;
			curr->next = curr->next->next;
			free(temp);
		}
		printf("\n");
	}
	return 0;
}
