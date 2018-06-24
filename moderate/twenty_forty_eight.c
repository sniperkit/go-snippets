#include <stdio.h>

static void slide(int *p, int n) {
	int i, j, c = 0, l = 0;
	for (i = 0; i < n; i++) {
		if (p[c] == 0) {
			for (j = c; j < n - 1; j++)
				p[j] = p[j + 1];
			p[n - 1] = 0;
		} else if (p[c] == l) {
			p[c - 1] <<= 1;
			for (j = c; j < n - 1; j++)
				p[j] = p[j + 1];
			p[n - 1] = 0;
			l = 0;
		} else {
			l = p[c++];
		}
	}
}

int main(int argc, char *argv[]) {
	FILE *fp;
	char d;
	int m[100], q[10], n, i, j;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((d = getc(fp)) != EOF) {
		switch (d) {
		case 'R':
			fseek(fp, 6, SEEK_CUR);
			break;
		case 'L':
		case 'D':
			fseek(fp, 5, SEEK_CUR);
			break;
		case 'U':
			fseek(fp, 3, SEEK_CUR);
		}
		fscanf(fp, "%d; ", &n);
		for (i = 0; i < n * n; i++)
			fscanf(fp, "%d%*c", &m[i]);
		for (i = 0; i < n; i++) {
			switch (d) {
			case 'L':
				for (j = 0; j < n; j++)
					q[j] = m[i * n + j];
				slide(q, n);
				for (j = 0; j < n; j++)
					m[i * n + j] = q[j];
				break;
			case 'R':
				for (j = 0; j < n; j++)
					q[n - 1 - j] = m[i * n + j];
				slide(q, n);
				for (j = 0; j < n; j++)
					m[i * n + j] = q[n - 1 - j];
				break;
			case 'U':
				for (j = 0; j < n; j++)
					q[j] = m[j * n + i];
				slide(q, n);
				for (j = 0; j < n; j++)
					m[j * n + i] = q[j];
				break;
			case 'D':
				for (j = 0; j < n; j++)
					q[n - 1 - j] = m[j * n + i];
				slide(q, n);
				for (j = 0; j < n; j++)
					m[j * n + i] = q[n - 1 - j];
			}
		}
		for (i = 0; i < n; i++) {
			for (j = 0; j < n; j++) {
				if (j > 0)
					putchar(' ');
				printf("%d", m[i * n + j]);
			}
			if (i < n - 1)
				putchar('|');
		}
		putchar('\n');
	}
	return 0;
}
