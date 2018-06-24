#include <stdio.h>
#include <sys/stat.h>

int main(int argc, char *argv[]) {
	struct stat st;
	int ret;

	if (argc != 2) {
		printf("Usage: %s [FILE]\n", argv[0]);
		return 1;
	}
	ret = stat(*++argv, &st);
	if (ret == 0)
		printf("%lld\n", (long long)st.st_size);
	return ret;
}
