#include <stdio.h>

static int robot(int f, int x, int y) {
	int ret = 0;
	if (x == 3 && y == 3)
		return 1;
	if (x > 0 && !(f & (1 << (x - 1 + 4 * y))))
		ret += robot(f | (1 << (x - 1 + 4 * y)), x - 1, y);
	if (y > 0 && !(f & (1 << (x + 4 * (y - 1)))))
		ret += robot(f | (1 << (x + 4 * (y - 1))), x, y - 1);
	if (x < 3 && !(f & (1 << (x + 1 + 4 * y))))
		ret += robot(f | (1 << (x + 1 + 4 * y)), x + 1, y);
	if (y < 3 && !(f & (1 << (x + 4 * (y + 1)))))
		ret += robot(f | (1 << (x + 4 * (y + 1))), x, y + 1);
	return ret;
}

int main(void) {
	printf("%d\n", robot(1, 0, 0));
	return 0;
}
