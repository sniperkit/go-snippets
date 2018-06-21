#include "list.h"
#include <stdlib.h>
#include <assert.h>
#include <string.h>
#include <assert.h>
#include <stdio.h>

#ifdef ___LIST_H___
struct List* ListCreate(COMPARE cmp, DESTROY dest)
{
	struct List* list = NULL;

	if (cmp == NULL || dest == NULL) 
		return NULL;

	list = (struct List*)malloc(sizeof (struct List));
	if (list == NULL) {
		return NULL;
	}

	memset(list, 0, sizeof(struct List));
	list->head = NULL;
	list->cmp = cmp;
	list->dest = dest;

	return list;
}

static inline int listIsEmpty(struct List *list) {
	assert(list != NULL);

	return list->head == NULL ? 1 : 0;
}

static inline struct Node* listGetLastNode(struct List* list) {
	struct Node* last = NULL;

	if (list == NULL) {
		assert(list != NULL);
		return NULL;
	}

	last = list->head;
	while (last->next != NULL) {
		last = last->next;
	}

	return last;
}

static inline struct Node* listGetPrevNode(struct List* list, void* key) {
	struct Node *prev = NULL;
	struct Node *temp = NULL;

	if (list == NULL || key == NULL || list->head == NULL) {
		assert(list != NULL && key != NULL && list->head != NULL);
		return NULL;
	}

	assert(ListIsNodeExist(list, key));

	prev = list->head;
	temp = prev;

	if (list->cmp(temp->data, key) == 0) {
		return NULL;
	} else {
		temp = temp->next;
	}

	while (temp != NULL) {
		if (list->cmp(temp->data, key) == 0) {
			break;
		}

		prev = prev->next;
		temp = prev->next;
	}

	if (temp == NULL) {
		return NULL;
	}

	return prev;
}

int ListIsNodeExist(struct List* list, void* key) {
	struct Node* temp = NULL;

	if (list == NULL || key == NULL) {
		assert(list != NULL && key != NULL);
		return -1;
	}

	if (listIsEmpty(list)) {
		return 0;
	}

	temp = list->head;
	while (temp != NULL) {
		if (list->cmp(temp->data, key) == 0) {
			return 1;
		}

		temp = temp->next;
	}

	return 0;
}

int ListInsertFront(struct List* list, void* data){
	struct Node *new = NULL;

	if (list == NULL || data == NULL) {
		assert(list != NULL && data != NULL);
		return -1;
	}

	new = (struct Node*)malloc(sizeof (struct Node));
	if (new == NULL) {
		assert(new != NULL);
		return -1;
	}

	memset(new, 0, sizeof(struct Node));
	new->data = data;
	new->next = NULL;

	new->next = list->head;
	list->head = new;

	return 0;
}

int ListInsertEnd(struct List* list, void *data){
	struct Node *new = NULL;
	struct Node *last = NULL;

	if (list == NULL || data == NULL) {
		assert(list != NULL && data != NULL);
		return -1;
	}

	new = (struct Node*) malloc(sizeof(struct Node));
	if (new == NULL) {
		assert(new != NULL);
		return -1;
	}
	
	memset(new, 0, sizeof(struct Node));
	new->next = NULL;
	new->data = data;

	if (listIsEmpty(list)) {
		list->head = new;
		return 0;
	}

	last = listGetLastNode(list);
	if (last == NULL) {
		assert(last != NULL);
		return -1;
	}

	last->next = new;
	return 0;
}

int ListInsertSorted(struct List* list, void* data){
	struct Node* new = NULL;
	struct Node* pos = NULL;
	struct Node* prev = NULL;

	if (list == NULL || data == NULL) {
		assert(list != NULL && data != NULL);
		return -1;
	}

	new = (struct Node*)malloc(sizeof(struct Node));
	if (new == NULL) {
		assert(new != NULL);
		return -1;
	}

	memset(new, 0, sizeof(struct Node));
	new->data = data;
	new->next = NULL;

	if(listIsEmpty(list)) {
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

	if (pos == NULL) {
		pos = listGetLastNode(list);
		pos->next = new;
		return 0;
	}

	prev = listGetPrevNode(list, pos->data);
	if (prev == NULL) {
		new->next = list->head;
		list->head = new;
	} else {
		new->next = prev->next;
		prev->next = new;
	}

	return 0;
}

int ListDeleteFirst(struct List* list){
	struct Node *del = NULL;
	if (list == NULL || list->head == NULL) {
		assert(list != NULL && list->head != NULL);
		return -1;
	}

	del = list->head;
	list->head = del->next;
	list->dest(del->data);
	free(del);

	return 0;
}

int ListDeleteLast(struct List *list){
	struct Node *del = NULL;
	struct Node *prev = NULL;

	if (list == NULL || list->head == NULL) {
		assert(list != NULL && list->head != NULL);
		return -1;
	}

	prev = list->head;
	del = prev->next;

	while (del != NULL && del->next != NULL) {
		prev = prev->next;
		del = prev->next;
	}

	if (del == NULL) {
		list->dest(prev->data);
		free(prev);
		return 0;
	}

	prev->next = NULL;
	list->dest(del->data);
	free(del);
	return 0;
}

int ListDeleteNode(struct List *list, void* key){
	struct Node* del = NULL;
	struct Node* prev = NULL;

	if (list == NULL || key == NULL || list->head == NULL) {
		assert(list != NULL && key != NULL && list->head != NULL);
		return -1;
	}

	if (!ListIsNodeExist(list, key)) {
		return -1;
	}

	prev = listGetPrevNode(list, key);
	if (prev == NULL) {
		return ListDeleteFirst(list);
	}

	del = prev->next;
	prev->next = prev->next->next;
	list->dest(del->data);
	free(del);

#if 0
	prev = list->head;
	del = prev;

	if (list->cmp(del->data, key) == 0) {
		return ListDeleteFirst(list);
	} else {
		del = del->next;
	}
	
	while(del != NULL) {
		if (list->cmp(del->data, key) == 0) {
			break;
		}

		prev = prev->next;
		del = prev->next;
	}

	if (del == NULL) {
		assert(del != NULL);
		return -1;
	}

	prev->next = del->next;
	list->dest(del->data);
	free(del);
#endif

	return 0;
}

int ListDeleteAll(struct List* list, void* key) {
	if (list == NULL || key == NULL) {
		assert(list != NULL && key != NULL);
		return -1;
	}

	if (listIsEmpty(list)) {
		return 0;
	}

	assert(ListIsNodeExist(list, key) == 1);

	while (ListIsNodeExist(list, key)) {
		ListDeleteNode(list, key);
	}

	return 0;
}

void* ListGet(struct List* list, void* key){
	struct Node* temp = NULL;

	if (list == NULL || key == NULL) {
		assert(list != NULL && key != NULL);
		return NULL;
	}

	if (listIsEmpty(list) || ListIsNodeExist(list, key) == 0) {
		return NULL;
	}

	temp = listGetPrevNode(list, key);
	if (temp == NULL) {
		return list->head->data;
	} 

	return temp->next->data;
}

int ListUpdate(struct List* list, void* key, void* data){
	struct Node* temp = NULL;

	if (list == NULL || key == NULL) {
		assert(list != NULL && key != NULL);
		return -1;
	}

	if (listIsEmpty(list) || ListIsNodeExist(list, key) == 0) {
		return -1;
	}

	temp = listGetPrevNode(list, key);
	if (temp == NULL) {
		list->head->data = data;
	}  else {
		temp->next->data = data;
	}

	return 0;
}

int ListForEach(struct List* list, OPERATE op){
	struct Node* temp = NULL;

	if (list == NULL || op == NULL) {
		assert(list != NULL && op != NULL);
		return -1;
	}

	temp = list->head;

	printf("=======================================\n");
	while(temp != NULL) {
		(*op)(temp->data);
		temp = temp->next;
	}

	return 0;
}

int ListReverse(struct List* list){
	struct Node* head = NULL;
	struct Node* next = NULL;
	struct Node* temp = NULL;
	if (list == NULL) {
		assert(list != NULL);
		return -1;
	}

	if (listIsEmpty(list)) {
		return 0;
	}

	head = list->head;
	next = head->next;
	head->next = NULL;

	//list has only one node.
	if (next == NULL) {
		return 0;
	}

	while (next != NULL) {
		temp = next;
		next = next->next;
		temp->next = head;
		head = temp;
	}

	list->head = head;
	return 0;
}

void ListFlush(struct List* list) {
	struct Node* temp = NULL;
	struct Node* next = NULL;

	if (list == NULL) {
		assert(list != NULL);
		return;
	}

	if (listIsEmpty(list)) {
		free(list);
		return;
	}

	temp = list->head;
	next = temp->next;
	while(temp != NULL) {
		list->dest(temp->data);
		free(temp);

		temp = next;
		if (next != NULL) {
			next = next->next;
		}
	}

	list->head = NULL;

	return;
}

void ListDestroy(struct List* list){
	ListFlush(list);
	free(list);
	return;
}

#endif /* ___LIST_H___ */
