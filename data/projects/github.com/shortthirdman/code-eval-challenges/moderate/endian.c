#include <stdio.h>

int main(void) {
	int i = 1;
	char *p = (char *)&i;

	puts(*p ? "LittleEndian" : "BigEndian");

	return 0;
}
