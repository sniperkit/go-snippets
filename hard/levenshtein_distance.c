#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

struct node {
	char		*data;
	unsigned long	length;
	struct node	*next;
};

static bool leveno(char *a, char *b) {
	bool d = false;
	while (*b != '\0') {
		if (*a != *b) {
			if (d)
				return false;
			d = true;
		} else {
			b++;
		}
		a++;
	}
	return true;
}

static bool levene(char *a, char *b) {
	bool d = false;
	while (*a != '\0') {
		if (*a != *b) {
			if (d)
				return false;
			d = true;
		}
		a++;
		b++;
	}
	return true;
}

int main(int argc, char *argv[]) {
	FILE *fp;
	char line[32];
	char root[6] = "hello";
	unsigned long netw = 0, strr = strlen(root), strl;
	struct node *look = NULL, *found = NULL, *curl = NULL, *curf = NULL;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fgets(line, 31, fp) != 0) {
		strl = strlen(line) - 1;
		struct node *temp = malloc(sizeof(struct node));
		char *w = malloc(strl);

		temp->data = strncpy(w, line, strl);
		temp->length = strl;
		temp->next = NULL;

		if ((strr == strl && levene(root, temp->data)) ||
		    (strr > strl && strr - strl < 2 && leveno(root, temp->data)) ||
		    (strr < strl && strl - strr < 2 && leveno(temp->data, root))) {
			if (curf)
				curf->next = temp;
			else
				found = temp;
			curf = temp;
		} else {
			if (curl)
				curl->next = temp;
			else
				look = temp;
			curl = temp;
		}
	}
	struct node *done = NULL, *prev = NULL;
	curf = found;
	while (curf) {
		curl = look;
		prev = NULL;
		while (curl) {
			if ((curf->length == curl->length && levene(curf->data, curl->data)) ||
			    (curf->length > curl->length && curf->length - curl->length < 2 && leveno(curf->data, curl->data)) ||
			    (curf->length < curl->length && curl->length - curf->length < 2 && leveno(curl->data, curf->data))) {
				if (prev) {
					prev->next = curl->next;
				} else {
					look = curl->next;
				}
				done = curl->next;
				curl->next = curf->next;
				curf->next = curl;
				curl = done;
			} else {
				prev = curl;
				curl = curl->next;
			}
		}
		done = curf;
		curf = curf->next;
		netw++;
		free(done->data);
		free(done);
	}
	printf("%lu\n", netw);
	return 0;
}
