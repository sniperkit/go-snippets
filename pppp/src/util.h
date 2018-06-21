#ifndef ___UTIL_H___
#define ___UTIL_H___
#include <stdio.h>
#include <assert.h>

typedef int (*COMPARE) (void* a, void* b);
typedef int (*DESTROY)(void* data);
typedef int (*OPERATE)(void* data);

#if defined(DEBUG)
#define WARN(C, format, ...) do { \
	if (!(C)) { \
		printf("[%s][%s][%d](%s): " format, __FILE__, __func__, __LINE__, #C, ##__VA_ARGS__);\
	}\
	assert(C);\
} while (0)

#define ASSERT(C, R, format, ...) do { \
	if (!(C)) { \
		printf("[%s][%s][%d](%s): " format, __FILE__, __func__, __LINE__, #C, ##__VA_ARGS__);\
		assert(C);\
		return R; \
	}\
} while (0)
#else
#define WARN(C, format, ...) do { \
	if (!(C)) { \
		printf("[%s][%s][%d](%s): " format, __FILE__, __func__, __LINE__, #C, ##__VA_ARGS__);\
	}\
} while (0)

#define ASSERT(C, R, format, ...) do { \
	if (!(C)) { \
		printf("[%s][%s][%d](%s): " format, __FILE__, __func__, __LINE__, #C, ##__VA_ARGS__);\
		return R; \
	}\
} while (0)
#endif  /* DEBUG */
#endif /* ___UTIL_H___ */
