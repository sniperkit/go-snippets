#include "vector.h"
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

#ifdef ___VECTOR_H___
struct Vector* VectorCreate(COMPARE cmp, DESTROY dest, uint32_t cap){
	struct Vector *new = NULL;
	ASSERT(cmp != NULL && dest != NULL && cap > 0, NULL, "%s\n", "Invalid input parameters");

	new = (struct Vector*)malloc(sizeof(struct Vector));
	ASSERT(new != NULL, NULL, "%s\n", "Out of memory");

	memset(new, 0, sizeof(struct Vector));
	new->data = malloc(cap * sizeof(size_t));
	if (new->data == NULL) {
		WARN(new->data != NULL,"%s\n", "Out of memory");
		free(new);
		return NULL;
	}

	memset(new->data, 0, cap);
	new->cap = cap;
	new->count = 0;
	new->cmp = cmp;
	new->dest = dest;
	new->size = sizeof(size_t);

	return new;
}

void VectorDestroy(struct Vector* vec){
	uint32_t i = 0;

	ASSERT(vec != NULL, , "%s\n", "Invalid input parameters");

	for (i = 0; i < vec->count; i++) {
		vec->dest(vec->data[i]);
	}

	free(vec->data);
	free(vec);
	return;
}

#define VECTOR_NEED_EXPAND(vec) ((vec)->count > (vec)->cap / 2)
#define VECTOR_NEED_COLLAPSE(vec) ((vec)->count < (vec)->cap / 2) 

static inline int vectorExpand(struct Vector* vec) {
	void *new = NULL;

	ASSERT(vec != NULL, -1, "%s\n", "Invalid input parameters");
	if (!VECTOR_NEED_EXPAND(vec)) {
		return 0;
	}

	new = malloc(vec->cap * 2 * vec->size);
	if (new == NULL) {
		WARN(new != NULL, "%s\n", "Cannot enlarge the vector due to out of memory");
		return -1;
	}

	memset(new, 0, vec->cap * 2 * vec->size);
	memcpy(new, vec->data, vec->count * vec->size);
	free(vec->data);
	vec->data = new;
	vec->cap = vec->cap * 2;

	return 0;
}

static inline int vectorCollapse(struct Vector* vec) {
	void* new = NULL;

	ASSERT(vec != NULL, -1, "%s\n", "Invalid input parameters");

	if (!VECTOR_NEED_COLLAPSE(vec)) {
		return 0;
	}

	new = malloc(vec->cap / 2 * vec->size);
	if (new == NULL) {
		WARN(new != NULL, "%s\n", "Cannot enlarge the vector due to out of memory");
		return -1;
	}

	memset(new, 0, vec->cap / 2 * vec->size);
	memcpy(new, vec->data, vec->count * vec->size);
	free(vec->data);
	vec->data = new;
	vec->cap = vec->cap / 2;

	return 0;
}

int VectorInsert(struct Vector* vec, void* data){
	ASSERT(vec != NULL && data != NULL, -1, "%s\n", "Invalid input parameters");

	return VectorInsertLast(vec, data);
}

int VectorInsertAt(struct Vector* vec, void* data, uint32_t idx) {
	int rv = -1;
	ASSERT(vec != NULL && data != NULL && idx > 0 && idx < vec->count, -1, "%s\n", "Invalid input parameters");

	if (VECTOR_NEED_EXPAND(vec)) {
		rv = vectorExpand(vec);
		ASSERT(rv == 0, -1, "%s\n", "Cannot expand the vector size");
	}

	memmove(&vec->data[idx+1], &vec->data[idx], (vec->count - idx) * vec->size); 
	vec->data[idx] = data;
	vec->count++;

	return 0;
}

int VectorInsertSorted(struct Vector* vec, void* data){
	int rv = -1;
	int idx = 0;
	ASSERT(vec != NULL && data != NULL, -1, "%s\n", "Invalid input parameters");
	ASSERT(vec->count < vec->cap, -1, "%s\n", "Vector is full, data overflow");

	if (VECTOR_NEED_EXPAND(vec)) {
		rv = vectorExpand(vec);
		ASSERT(rv == 0, -1, "%s\n", "Cannot expand the vector size");
	}

	return VectorInsertLast(vec, data);
}

int VectorInsertFront(struct Vector* vec, void* data){
	int rv = -1;
	ASSERT(vec != NULL && data != NULL, -1, "%s\n", "Invalid input parameters");
	ASSERT(vec->count < vec->cap, -1, "%s\n", "Vector is full, data overflow");

	if (VECTOR_NEED_EXPAND(vec)) {
		rv = vectorExpand(vec);
		ASSERT(rv == 0, -1, "%s\n", "Cannot expand the vector size");
	}

	memmove(&vec->data[1], vec->data, vec->count * vec->size);
	vec->data[0] = data;
	vec->count++;
	return 0;
}

int VectorInsertLast(struct Vector* vec, void* data){
	int rv = -1;
	ASSERT(vec != NULL && data != NULL, -1, "%s\n", "Invalid input parameters");
	ASSERT(vec->count < vec->cap, -1, "%s\n", "Vector is full, data overflow");

	if (VECTOR_NEED_EXPAND(vec)) {
		rv = vectorExpand(vec);
		ASSERT(rv == 0, -1, "%s\n", "Cannot expand the vector size");
	}

	vec->data[vec->count] = data;
	vec->count++;

	return 0;
}

int VectorDeleteFirst(struct Vector* vec){
	ASSERT(vec != NULL, -1, "%s\n", "Invalid input parameters");
	ASSERT(vec->count < vec->cap, -1, "%s\n", "Vector state is invalid, data overflow may be happend");

	memmove(vec->data, &vec->data[1], (vec->count - 1) * vec->size);
	vec->count--;

	return vectorCollapse(vec);
}

int VectorDeleteLast(struct Vector* vec){
	ASSERT(vec != NULL, -1, "%s\n", "Invalid input parameters");
	ASSERT(vec->count < vec->cap, -1, "%s\n", "Vector state is invalid, data overflow may be happend");
	vec->count--;
	return vectorCollapse(vec);
}

static inline int vectorIsNodeExist(struct Vector* vec, void* key) {
	uint32_t i = 0;
	ASSERT(vec != NULL && key != NULL, -1, "%s\n", "Invalid input parameters");
	ASSERT(vec->count < vec->cap, -1, "%s\n", "Vector state is invalid, data overflow may be happend");

	for(i = 0; i < vec->count; i++) {
		if (vec->cmp(vec->data[i], key) == 0) {
			return 1;
		}
	}

	return 0;
}

int VectorDeleteAt(struct Vector* vec, uint32_t idx) {
	ASSERT(vec != NULL && idx <= vec->count, -1, "%s\n", "Invalid input parameters");
	ASSERT(vec->count < vec->cap, -1, "%s\n", "Vector state is invalid, data overflow may be happend");

	if (idx == 0) {
		return	VectorDeleteFirst(vec);
	} 
	
	if (idx == vec->count - 1) {
		return VectorDeleteLast(vec);
	} 

	memmove(&vec->data[idx], &vec->data[idx+1], (vec->count - idx - 1) * vec->size);
	vec->count--;
	return vectorCollapse(vec);
}

int VectorDelete(struct Vector* vec, void* key){
	uint32_t i = 0;
	ASSERT(vec != NULL && key != NULL, -1, "%s\n", "Invalid input parameters");
	ASSERT(vec->count < vec->cap, -1, "%s\n", "Vector state is invalid, data overflow may be happend");


	if (!vectorIsNodeExist(vec, key)) {
		WARN(0 == 1, "%s\n", "Try to delete an unexisted entry");
		return -1;
	}

	for(i = 0; i < vec->count; i++) {
		if (vec->cmp(vec->data[i], key) == 0) {
			break;
		}
	}

	if (i >= vec->count) {
		WARN(0 == 1, "%s\n", "Try to delete an unexisted entry");
		return -1;
	}	

	return VectorDeleteAt(vec, i);
}

static inline int vectorIsEmpty(struct Vector *vec) {
	ASSERT(vec != NULL, -1, "%s\n", "Invalid input parameters");
	return vec->count == 0 ? 1 : 0;
}

void* VectorGetFirst(struct Vector* vec){
	ASSERT(vec != NULL, NULL, "%s\n", "Invalid input parameters");
	ASSERT(!vectorIsEmpty(vec), NULL, "%s\n", "Try to Get from an empty vector");

	return vec->data;
}

void* VectorGetLast(struct Vector* vec){
	ASSERT(vec != NULL, NULL, "%s\n", "Invalid input parameters");
	ASSERT(!vectorIsEmpty(vec), NULL, "%s\n", "Try to Get from an empty vector");
	return vec->data[(vec->count-1)];
}

void* VectorGet(struct Vector* vec, void* key){
	uint32_t i = 0;
	ASSERT(vec != NULL && key != NULL, NULL, "%s\n", "Invalid input parameters");
	ASSERT(!vectorIsEmpty(vec), NULL, "%s\n", "Try to Get from an empty vector");

	for(i = 0; i < vec->count; i++) {
		if (vec->cmp(vec->data[i], key) == 0) {
			break;
		}
	}

	if (i >= vec->count) {
		WARN(0 == 1, "%s\n", "Entry not exist");
		return NULL;
	}	

	return VectorGetAt(vec, i);
}

void *VectorGetAt(struct Vector* vec, uint32_t idx) {
	ASSERT(vec != NULL && idx <= vec->count, NULL, "%s\n", "Invalid input parameters");
	return vec->data[idx];
}

void VectorForEach(struct Vector* vec, OPERATE op){
	uint32_t idx = 0;
	ASSERT(vec != NULL && op != NULL, , "%s\n", "Invalid input parameters");
	
	printf("==============Dump Vector===vec->size: %d, vec->count: %d, vec->cap: %d========================\n", vec->size, vec->count, vec->cap);
	for (idx = 0; idx < vec->count; idx++) {
		(*op)(vec->data[idx]);
	}

	return ;
}

void VectorSort(struct Vector* vec){
	ASSERT(vec != NULL, , "%s\n", "Invalid input parameters");
	return ;
}
#endif /* ___VECTOER_H___ */

