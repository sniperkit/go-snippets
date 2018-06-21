#ifndef __prompt_H__
#define __prompt_H__
void prompt(void);

struct queue 
{
    int size;
    int head;
    int tail;
    int count;
    int* data;
};

struct queue* create_queue(int size);
void destroy_queue(struct queue* q);
int enqueue(struct queue* q, int data);
int dequeue(struct queue* q, int* data);
int queue_foreach(struct queue* q, void (*fn)(int data));
#endif /* __prompt_H__ */
