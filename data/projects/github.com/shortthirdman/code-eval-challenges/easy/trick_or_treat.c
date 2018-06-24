#include <stdio.h>

int main(int argc, char *argv[]) {
	FILE *fp;
	unsigned v, z, w, h;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	fp = fopen(*++argv, "r");
	while (fscanf(fp,
		      "Vampires: %d, Zombies: %d, Witches: %d, Houses: %d\n",
		      &v, &z, &w, &h) != EOF)
		printf("%d\n", v + z + w == 0 ? 0 : (v * 3 + z * 4 + w * 5) * h / (v + z + w));
	return 0;
}
