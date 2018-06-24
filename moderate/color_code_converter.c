#include <math.h>
#include <stdio.h>
#include <stdlib.h>

static float mod2(float a) {
	return a - (float)((int)(a / 2) * 2);
}

static unsigned rnint(float a) {
	return (unsigned)(a * 255 + 0.5);
}

int main(int argc, char *argv[]) {
	FILE *fp;
	char ch;
	unsigned r, g, b;
	float h, s, l, v, c, m, y, k;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while ((ch = getc(fp)) != EOF) {
		r = 0;
		g = 0;
		b = 0;
		switch (ch) {
		case 'H':
			fseek(fp, 1, SEEK_CUR);
			ch = getc(fp);
			if (ch == 'L') {
				fscanf(fp, "(%f,%f,%f)%*c", &h, &s, &l);
				s /= 100;
				l /= 100;
				c = (1 - fabsf(2 * l - 1)) * s;
				m = l - c / 2;
			} else {
				fscanf(fp, "(%f,%f,%f)%*c", &h, &s, &v);
				s /= 100;
				v /= 100;
				c = v * s;
				m = v - c;
			}
			h /= 60;
			y = c * (1 - fabsf(mod2(h) - 1));
			if (h < 1) {
				r = rnint(c + m);
				g = rnint(y + m);
				b = rnint(m);
			} else if (h < 2) {
				r = rnint(y + m);
				g = rnint(c + m);
				b = rnint(m);
			} else if (h < 3) {
				r = rnint(m);
				g = rnint(c + m);
				b = rnint(y + m);
			} else if (h < 4) {
				r = rnint(m);
				g = rnint(y + m);
				b = rnint(c + m);
			} else if (h < 5) {
				r = rnint(y + m);
				g = rnint(m);
				b = rnint(c + m);
			} else if (h < 6) {
				r = rnint(c + m);
				g = rnint(m);
				b = rnint(y + m);
			}
			break;
		case '(':
			fscanf(fp, "%f,%f,%f,%f)%*c", &c, &m, &y, &k);
			r = rnint((1 - c) * (1 - k));
			g = rnint((1 - m) * (1 - k));
			b = rnint((1 - y) * (1 - k));
			break;
		case '#':
			fscanf(fp, "%02x%02x%02x%*c", &r, &g, &b);
		}
		printf("RGB(%u,%u,%u)\n", r, g, b);
	}
	return 0;
}
