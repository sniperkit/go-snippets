#ifndef ___VECTOR_H___
#define ___VECTOR_H___

#include<inttypes.h>
#include "util.h"

struct Vector {
	uint32_t cap;
	uint32_t count;
	uint32_t size;
	COMPARE cmp;
	DESTROY dest;
	void **data;
};

struct Vector* VectorCreate(COMPARE cmp, DESTROY dest, uint32_t cap);
void VectorDestroy(struct Vector* vec);
int VectorInsert(struct Vector* vec, void* data);
int VectorInsertSorted(struct Vector* vec, void* data);
int VectorInsertFront(struct Vector* vec, void* data);
int VectorInsertLast(struct Vector* vec, void* data);
int VectorInsertAt(struct Vector* vec, void* data, uint32_t idx);
int VectorDeleteFirst(struct Vector* vec);
int VectorDeleteLast(struct Vector* vec);
int VectorDelete(struct Vector* vec, void* key);
int VectorDeleteAt(struct Vector* vec, uint32_t idx);
void* VectorGetFirst(struct Vector* vec);
void* VectorGetLast(struct Vector* vec);
void* VectorGet(struct Vector* vec, void* key);
void* VectorGetAt(struct Vector* vec, uint32_t idx);
void VectorForEach(struct Vector* vec, OPERATE op);
void VectorSort(struct Vector* vec);
#endif /* ___VECTOR_H___ */
