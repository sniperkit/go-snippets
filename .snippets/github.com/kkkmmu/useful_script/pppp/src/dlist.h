#ifndef ___DLIST_H___
#define ___DLIST_H___

#include "util.h"
struct DList {
	struct  DNode* head;
	DESTROY dest;
	COMPARE cmp;
};

struct DNode {
	struct DNode* prev;
	struct DNode* next;
	void* data;
};

struct DList* DListCreate(COMPARE cmp, DESTROY dest);
void DListDestroy(struct DList* list);
int DListInsertFront(struct DList* list, void* data);
int DListInsertEnd(struct DList* list, void* data);
int DListInsert(struct DList* list, void* data);
int DListInsertSorted(struct DList* list, void* data);
int DListDeleteFirst(struct DList* list);
int DListDeleteLast(struct DList* list);
int DListDelete(struct DList* list, void* key);
void* DListGetFirst(struct DList* list);
void* DListGetLast(struct DList* list);
void* DListGet(struct DList* list, void* key);
int DListSort(struct DList* list);
int DListForEach(struct DList* list, OPERATE op);
void DListReverse(struct DList* list);
#endif /* ___DLIST_H___ */
