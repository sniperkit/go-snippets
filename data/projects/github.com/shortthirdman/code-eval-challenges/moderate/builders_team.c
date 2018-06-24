#include <stdio.h>

static const unsigned hor[30] = { 0x42, 0x84, 0x108, 0x210,
				  0x840, 0x1080, 0x2100, 0x4200,
				  0x10800, 0x21000, 0x42000, 0x84000,
				  0x210000, 0x420000, 0x840000, 0x1080000,
				  0x1806, 0x300c, 0x6018,
				  0x300c0, 0x60180, 0xc0300,
				  0x601800, 0xc03000, 0x1806000,
				  0x7000e, 0xe001c,
				  0xe001c0, 0x1c00380,
				  0x1e0001e };
static const unsigned ver[30] = { 0x6, 0xc, 0x18, 0x30,
				  0xc0, 0x180, 0x300, 0x600,
				  0x1800, 0x3000, 0x6000, 0xc000,
				  0x30000, 0x60000, 0xc0000, 0x180000,
				  0x14a, 0x294, 0x528,
				  0x2940, 0x5280, 0xa500,
				  0x52800, 0xa5000, 0x14a000,
				  0x4a52, 0x94a4,
				  0x94a40, 0x129480,
				  0x118c62 };

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned a, b, t;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp, "%d %d", &a, &b) != EOF) {
		unsigned h = 0, v = 0, r = 0, i;
		do {
			if (a > b) {
				t = a;
				a = b;
				b = t;
			}
			if (a + 1 == b)
				h |= 1 << a;
			else
				v |= 1 << a;
		} while (getc(fp) != '\n' && fscanf(fp, "| %d %d", &a, &b) != EOF);
		for (i = 0; i < 30; i++)
			if ((h & hor[i]) == hor[i] && (v & ver[i]) == ver[i])
				r++;
		printf("%d\n", r);
	}
	return 0;
}
