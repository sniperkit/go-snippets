#ifndef ___LIST_H___
#define ___LIST_H___

#include "util.h"

struct Node {
	struct Node* next;
	void *data;
};

struct List {
	struct Node* head;
	COMPARE cmp;
	DESTROY dest;
};

struct List* ListCreate(COMPARE cmp, DESTROY dest);
int ListInsertFront(struct List* list, void* data);
int ListInsertEnd(struct List* list, void* data);
int ListInsertSorted(struct List* list, void* data);
int ListDeleteFirst(struct List* list);
int ListDeleteLast(struct List* list);
int ListDeleteNode(struct List* list, void* key);
int ListDeleteAll(struct List* list, void* key);
int ListIsNodeExist(struct List* list, void* key);
void* ListGet(struct List* list, void* key);
int ListUpdate(struct List* list, void* key, void* data);
int ListForEach(struct List* list, OPERATE op);
int ListReverse(struct List* list);
void ListFlush(struct List* list);
void ListDestroy(struct List* list);
#endif /* ___LIST_H___ */
