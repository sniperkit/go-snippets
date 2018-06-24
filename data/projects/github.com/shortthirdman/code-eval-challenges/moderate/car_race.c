#include <stdio.h>
#include <stdlib.h>

struct segment {
	float length;
	float angle;
};

struct car {
	unsigned id;
	float    time;
};

static int cmpcar(const void *p1, const void *p2) {
	struct car *c1 = (struct car *)p1;
	struct car *c2 = (struct car *)p2;
	return (int)(100 * c1->time - 100 * c2->time);
}

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned i, s = 0, sbs = 16, n = 0, id;
	float vmax, accel, brake;
	char c;
	struct segment *sb = malloc(sbs * sizeof(struct segment));
	struct car cars[40];

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	do {
		if (s == sbs) {
			sbs += sbs / 2;
			sb = realloc(sb, sbs * sizeof(struct segment));
		}
		fscanf(fp, "%f %f", &sb[s].length, &sb[s].angle);
		s++;
	} while ((c = getc(fp)) == ' ');

	while (fscanf(fp, "%d %f %f %f", &id, &vmax, &accel, &brake) != EOF) {
		float laptime = 0, vinit = 0;
		for (i = 0; i < s; i++) {
			float len = sb[i].length, away, vend, bway, angle = sb[i].angle;
			away = accel * (1 - (vinit / vmax)) * (vinit + vmax) / (2 * 3600);
			vend = vmax * (1 - (angle / 180));
			bway = brake * (1 - (vend / vmax)) * (vend + vmax) / (2 * 3600);
			laptime += 3600 * (2 * away / (vinit + vmax) + (len - away - bway) / vmax +
					   2 * bway / (vend + vmax));
			vinit = vend;
		}
		cars[n].id = id;
		cars[n++].time = laptime;
	}
	qsort(cars, n, sizeof(struct car), cmpcar);
	for (i = 0; i < n; i++) {
		printf("%d %.2f\n", cars[i].id, cars[i].time);
	}
	free(sb);
	return 0;
}
