#include <stdio.h>


int main() {
outer:
	for (int i = 1; i < 100; i++) {
		if (i % 2 == 0)
			goto outer;

		printf("%d\n", i);
	}
	return 0;
}
