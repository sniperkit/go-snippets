#include "dlist.h"
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#if defined(___DLIST_H___)
struct DList* DListCreate(COMPARE cmp, DESTROY dest){
	struct DList* new = NULL;

	ASSERT(cmp != NULL && dest != NULL, NULL, "%s\n", "cmp and dest is necessary!");

	new = (struct DList*)malloc(sizeof(struct DList));
	ASSERT(new != NULL, NULL, "%s\n", "Cannot alloc memory for new list");

	memset(new, 0, sizeof(struct DList));
	new->cmp = cmp;
	new->dest = dest;
	new->head = NULL;

	return new;
}

void DListDestroy(struct DList* list){
	struct DNode* temp = NULL;
	struct DNode* next = NULL;

	ASSERT(list != NULL, , "%s\n", "Try to destroy an invalid list");

	temp = list->head;
	while(temp != NULL) {
		next = temp->next;
		list->dest(temp->data);
		free(temp);
		temp = NULL;
		if (next != NULL) {
			temp = next;
		}
	}

	list->head = NULL;
	free(list);
	return;
}

int DListInsertFront(struct DList* list, void* data){
	struct DNode* new = NULL;

	ASSERT(list != NULL && data != NULL, -1, "%s\n", "list and data should not be NULL");

	new = (struct DNode*)malloc(sizeof(struct DNode));
	ASSERT(new != NULL, -1, "%s\n", "Cannot alloc mem for new node");
	memset(new, 0, sizeof(struct DNode));
	new->data = data;
	new->next = NULL;
	new->prev = NULL;

	new->next = list->head;
	if (list->head != NULL) {
		list->head->prev = new;
	}
	list->head = new;

	return 0;
}

static inline int dlistIsEmpty(struct DList* list) {
	ASSERT(list != NULL, -1, "%s\n", "Invalid list");

	return list->head == NULL ? 1 : 0;
}

static inline struct DNode* dlistGetLastNode(struct DList* list) {
	struct DNode* temp = NULL;

	ASSERT(list != NULL, NULL, "%s\n", "Invalid list");

	if (dlistIsEmpty(list)) {
		return NULL;
	}

	temp = list->head;
	while (temp->next != NULL) {
		temp = temp->next;
	}

	return temp;
}

int DListInsertEnd(struct DList* list, void* data){
	struct DNode* new = NULL;
	struct DNode* pos = NULL;

	ASSERT(list != NULL && data != NULL, -1, "%s\n", "list and data should not be NULL");
	
	new = (struct DNode*)malloc(sizeof(struct DNode));
	ASSERT(list != NULL, -1, "%s\n", "Cannot alloc mem for new node");
	memset(new, 0, sizeof(struct DNode));
	new->next = NULL;
	new->prev = NULL;
	new->data = data;

	pos = dlistGetLastNode(list);
	if (pos == NULL) {
		list->head = new;
	} else {
		pos->next = new;
		new->prev = pos;
	}

	return 0;
}

int DListInsert(struct DList* list, void* data){
	return DListInsertEnd(list, data);
}

int DListInsertSorted(struct DList* list, void* data){
	struct DNode* pos = NULL;
	struct DNode* new = NULL;

	ASSERT(list != NULL, -1, "%s\n", "Inavlid list to insert entry.");

	new = (struct DNode*)malloc(sizeof(struct DList));	
	ASSERT(new != NULL, -1, "%s\n", "Cannot create new node for out of memory");

	memset(new, 0, sizeof(struct DNode));
	new->next = NULL;
	new->prev = NULL;
	new->data = data;

	if (dlistIsEmpty(list)) {
		list->head = new;
		return 0;
	}

	pos = list->head;
	while(pos != NULL) {
		if (list->cmp(new->data, pos->data) > 0) {
			pos = pos->next;
			continue;
		}

		break;
	}

	//1. At last.
	if (pos == NULL) {
		pos = dlistGetLastNode(list);
		pos->next = new;
		new->prev = pos;
		return 0;
	}

	//2. Insert At first.
	if (pos->prev == NULL) {
		list->head->prev = new;
		new->next = list->head;
		list->head = new;
		return 0;
	}

	//3.
	new->next = pos->prev->next;
	new->prev = pos->prev;
	new->next->prev = new;
	new->prev->next = new;
	return 0;
}

int DListDeleteFirst(struct DList* list){
	struct DNode* del = NULL;
	struct DNode* temp = NULL;

	ASSERT(list != NULL, -1, "%s\n", "Invalid List to delete");

	if (dlistIsEmpty(list)) {
		WARN(list->head != NULL, "%s\n", "Try to delete from an empty list");
		return 0;
	}

	del = list->head;
	temp = del->next;
	if (temp != NULL) {
		temp->prev = NULL;
	}
	list->head = temp;

	list->dest(del->data);
	free(del);
	return 0;
}

int DListDeleteLast(struct DList* list){
	struct DNode* del = NULL;

	ASSERT(list != NULL, -1, "%s\n", "Invalid input list");

	if (dlistIsEmpty(list)) {
		WARN(list->head != NULL, "%s\n", "Try to delete entry from empty list");
		return 0;
	}

	del = dlistGetLastNode(list);
	if (del->prev != NULL) {
		del->prev->next = NULL;
	} else {
		list->head = NULL;
	}

	list->dest(del->data);
	free(del);
		
	return 0;
}

static inline int dlistIsNodeExist(struct DList* list, void* key) {
	struct DNode* temp = NULL;
	
	ASSERT(list != NULL && key != NULL, -1, "%s\n", "Invalid input parameters");
	temp = list->head;
	while (temp != NULL) {
		if (list->cmp(temp->data, key) == 0) {
			return 1;
		}
		temp = temp->next;
	}

	return 0;
}

struct DNode* DListGetNode(struct DList* list, void* key) {
	struct DNode* temp = NULL;
	
	ASSERT(list != NULL && key != NULL, NULL, "%s\n", "Invalid input parameters");
	temp = list->head;
	while (temp != NULL) {
		if (list->cmp(temp->data, key) == 0) {
			return temp;
		}
		temp = temp->next;
	}

	return NULL;
}

int DListDelete(struct DList* list, void* key){
	struct DNode* del = NULL;
	struct DNode* temp = NULL;

	ASSERT(list != NULL && key != NULL, -1, "%s\n", "Invalid input parameters");

	if (!dlistIsNodeExist(list, key)) {
		WARN(0 == 1, "%s\n", "Try to delete an unexisted entry");
		return -1;
	}

	del = DListGetNode(list, key);
	if (del == NULL) {
		WARN(0 == 1, "%s\n", "Try to delete an unexisted entry");
		return -1;
	}

	if (del->prev == NULL && del->next == NULL) {
		list->head = NULL;
	} else if (del->prev == NULL) {
		list->head  = del->next;
		del->next->prev = NULL;
	} else if (del->next == NULL) {
		del->prev->next = NULL;
	} else {
		temp = del->next;
		temp->prev = del->prev;
		temp->prev->next = temp;
	}

	list->dest(del->data);
	free(del);

	return 0;
}

void* DListGetFirst(struct DList* list){
	ASSERT(list != NULL, "NULL", "%s\n", "Invalid input paramters");

	if (dlistIsEmpty(list)) {
		WARN(1 == 0, "%s\n", "Try to get from empty list");
		return NULL;
	}

	return list->head->data;
}

void* DListGetLast(struct DList* list){
	struct DNode* last = NULL;

	ASSERT(list != NULL, "NULL", "%s\n", "Invalid input paramters");

	if (dlistIsEmpty(list)) {
		WARN(1 == 0, "%s\n", "Try to get from empty list");
		return NULL;
	}

	last = dlistGetLastNode(list);
	if (last == NULL) {
		WARN(1 == 0, "%s\n", "Try to get from empty list");
		return NULL;
	}

	return last->data;
}

void* DListGet(struct DList* list, void* key){
	struct DNode* pos = NULL;

	ASSERT(list != NULL, NULL, "%s\n", "Invalid input paramters");

	if (dlistIsEmpty(list)) {
		WARN(1 == 0, "%s\n", "Try to get from empty list");
		return NULL;
	}

	pos = DListGetNode(list, key);
	if (pos == NULL) {
		WARN(1 == 0, "%s\n", "Cannot find a node by this key");
		return NULL;
	}

	return pos->data;
}

int DListForEach(struct DList* list, OPERATE op){
	struct DNode* temp = NULL;

	ASSERT(list != NULL, 0, "%s\n", "Invalid input paramters");

	if (dlistIsEmpty(list)) {
		WARN(1 == 0, "%s\n", "Try to traverse an empty list");
		return -1;
	}

	printf("=======================================\n");
	temp = list->head;
	while (temp != NULL) {
		(*op)(temp->data);
		temp = temp->next;
	}

	return 0;
}

void DListReverse(struct DList* list){
	struct DNode* temp = NULL;
	struct DNode* next = NULL;
	struct DNode* head = NULL;

	ASSERT(list != NULL, , "%s\n", "Invalid input paramters");
	
	if (dlistIsEmpty(list)) {
		return;
	}

	head = list->head;
	next = head->next;
	head->next = NULL;
	while (next != NULL) {
		temp = next;
		next = next->next;
		temp->next = head;
		head->prev = temp;
		head = temp;
	}

	list->head = head;

	return;
}

#endif /* ___DLIST_H___ */
