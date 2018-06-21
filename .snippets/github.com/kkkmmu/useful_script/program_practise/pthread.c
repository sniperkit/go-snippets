#include <pthread.h>
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <errno.h>
#include <ctype.h>
#include <llog.h>
#include "test.h"
#include "random.h"
#include <signal.h>
#include <time.h>
#include <pthread.h>
#include <assert.h>

#define handle_error_en(en, msg) \
    do { errno = en; perror(msg); exit(EXIT_FAILURE); } while (0)

#define handle_error(msg) \
    do { perror(msg); exit(EXIT_FAILURE); } while (0)

#define QSIZE 100000

pthread_mutex_t qlock = PTHREAD_MUTEX_INITIALIZER;

struct queue* Q = NULL;

struct queue* create_queue(int size) 
{
    struct queue* new = NULL;
    new = (struct queue*)malloc(sizeof(*new));
    if (new == NULL)
        return NULL;

    new->data = (int *)malloc(size*sizeof(int));
    if (new->data == NULL) {
        free(new);
        return NULL;
    }

    memset(new->data, 0, size);
    new->head = 0;
    new->tail = 0;
    new->size = size;
    new->count = 0;

    return new;
}

void destroy_queue(struct queue* q)
{
    if (q == NULL)
        return;

    //pthread_mutex_lock(&qlock);
    if (q->data)
        free(q->data);

    free(q);
    //pthread_mutex_unlock(&qlock);
} 

static int isfull(struct queue* q)
{
    int ret = 0;
    if (q == NULL)
        return 1;
    assert(q->count <= QSIZE);

    //pthread_mutex_lock(&qlock);
    ret = q->count == q->size ? 1 : 0;
    //pthread_mutex_unlock(&qlock);
    return ret;
}

static int isempty(struct queue* q)
{
    int ret = 0;
    if (q == NULL)
        return 0;

    assert(q->count >= 0);
    //pthread_mutex_lock(&qlock);
    ret = q->count == 0 ? 1 : 0;
    //pthread_mutex_unlock(&qlock);
    return ret;
}

int enqueue(struct queue* q, int data)
{
    //pthread_mutex_lock(&qlock);
    if (isfull(q))
    {
        //pthread_mutex_unlock(&qlock);
        return -1;
    }
    q->data[(q->tail)%q->size] = data;
    q->tail = (q->tail+1)%q->size;
    q->count++;
    //pthread_mutex_unlock(&qlock);
    llog(INFO, "+++(%d) is enqueed!(%d) by: %u\n", data, q->count, pthread_self());
    return 0;
}

int dequeue(struct queue* q, int* data)
{
    //pthread_mutex_lock(&qlock);
    if (isempty(q))
    {
     //   pthread_mutex_unlock(&qlock);
        return  -1;
    }
    *data = (q->data[q->head]);
    q->head = (q->head+1)%q->size;

    q->count--;
    //pthread_mutex_unlock(&qlock);
    llog(INFO, "---(%d) is dequeued! (%d) by: %u\n", *data, q->count, pthread_self());
    return 0;
}

int queue_foreach(struct queue* q, void (*fn)(int data))
{
    if (q == NULL || q->count == 0)
        return 0;

    printf("============================\n");
    printf("queue->head: %d\n", q->head);
    printf("queue->tail: %d\n", q->tail);
    printf("queue->size: %d\n", q->size);
    printf("queue->count: %d\n", q->count);
    printf("============================\n");
    int i = 0;
    for (; i < q->count; i++)
    {
        (*fn)((q->data[(i + q->head)% q->size]));
    }
    return 0;
}

//int pthread_create(pthread_t *thread, const pthread_attr_t *attr,
 //       void *(*start_routine) (void *), void *arg);


#define CLOCKID CLOCK_REALTIME

void produce_timer(union sigval v)  
{  
    pthread_mutex_lock(&qlock);
    enqueue(Q, random_uint32());
    pthread_mutex_unlock(&qlock);
} 

void* produce(void* data)
{
#if 0
    timer_t timerid;
    struct sigevent evp;
    memset(&evp, 0, sizeof(struct sigevent));
    evp.sigev_value.sival_int = 120;
    evp.sigev_notify = SIGEV_THREAD;
    evp.sigev_notify_function = produce_timer; 

    timer_create(CLOCKID, &evp, &timerid); 
    struct itimerspec it;  
    it.it_interval.tv_sec = 1;  
    it.it_interval.tv_nsec = 0;  
    it.it_value.tv_sec = 1;  
    it.it_value.tv_nsec = 0;  

    if (timer_settime(timerid, 0, &it, NULL) == -1)  
    {  
        perror("fail to timer_settime");  
        exit(-1);  
    }  

    pause();
#else
    long i = 1000000000;
    long j = 0;

    for (; j < i; j++) {
        //pthread_mutex_lock(&qlock);
        enqueue(Q, random_uint32());
        //pthread_mutex_unlock(&qlock);
        usleep(20000);
    }
#endif
    return NULL;
}

void consume_timer(union sigval v)
{
    int p;
    int ret =0;
    ret = pthread_mutex_lock(&qlock);
    if (ret != 0)
    {
        llog(ERROR, "Error happened when lock: %d\n", ret);
        return ;
    }

    dequeue(Q, &p);
    pthread_mutex_unlock(&qlock);
    if (ret != 0)
    {
        llog(ERROR, "Error happened when unlock: %d\n", ret);
        return ;
    }
}

void* consume(void* data)
{
#if 0
    timer_t timerid;
    struct sigevent evp;
    memset(&evp, 0, sizeof(struct sigevent));
    evp.sigev_value.sival_int = 120;
    evp.sigev_notify = SIGEV_THREAD;
    evp.sigev_notify_function = consume_timer; 

    timer_create(CLOCKID, &evp, &timerid); 
    struct itimerspec it;  
    it.it_interval.tv_sec = 1;  
    it.it_interval.tv_nsec = 0;  
    it.it_value.tv_sec = 1;  
    it.it_value.tv_nsec = 0;  

    if (timer_settime(timerid, 0, &it, NULL) == -1)  
    {  
        perror("fail to timer_settime");  
        exit(-1);  
    }  

    pause();
#else
    long i = 100000000000;
    long j = 0;
    int k = 0;

    for (; j < i; j++) {

        //pthread_mutex_lock(&qlock);
        dequeue(Q, &k);
        //pthread_mutex_unlock(&qlock);
        usleep(20000);
    }

#endif
    return NULL;
}

pthread_t start_producer()
{
    pthread_t p;
    pthread_create(&p, NULL,  produce, NULL);
    return p;
}

pthread_t start_consumer()
{
    pthread_t p;
    pthread_create(&p, NULL,  consume, NULL);
    return p;
}

void print(int data) 
{
    printf("%d\n", data);
}

void test(void) 
{
    long p = 10000;
    long c = 1000;
    long j = 0;
    random_init();
    Q = create_queue(QSIZE);
    for(; j < p; j++)
        start_producer();
    //pthread_join(p, NULL);
    
    for(j = 0; j < c; j++)
        start_consumer();
  
    pause();
}
