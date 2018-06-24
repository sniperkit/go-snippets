//
// Created by yang on 24/05/2017.
//
#include <iostream>
#include "stdio.h"
#include "stdlib.h"

#define MAXSIZE 14
#define n 3
#define fee 10
using namespace std;

struct car {
    char state;
    int num;
    int time;
};

struct ParkStack {
    struct car G[n];
    int top;
};

struct rangweicar {
    int num;
    int time;
};

struct SqStack {
    struct rangweicar H[MAXSIZE];
    int s_top;
};

struct QNODE {
    int data;
    QNODE *next;
};

struct LinkQueue {
    QNODE *front,*rear;
    int count;
};

void Append_cars(ParkStack *s, LinkQueue *q, car a) {
    QNODE *t;
    if (s->top != n-1) {
        (s->top)++;
        (s->G[s->top]).state = a.state;
        (s->G[s->top]).num = a.num;
        (s->G[s->top]).time = a.time;
    }
    else {
        cout << "停车场已满!\n";
        t = new (QNODE);
        t->data = a.num;
        t-> next = NULL;
        q->rear->next = t;
        q->rear = t;
        q->count++;
    }
}

int D_cars (ParkStack *s, LinkQueue *q, car d) {
    int i, j, l, x, y;
    QNODE *p;
    SqStack *k;
    if (d.num == (s->G[s->top]).num) {
        x = d.time - (s->G[s->top]).time;
        y = fee * x;
        cout<<"停车时间为："<<x<<"小时，停车费用为："<<y<<"元!\n";
        if (q->count == 0) {
            cout<<"便道为空！"<<endl;
            s->top = s->top - 1;
            return 0;
        }
        else {
            p = q->front->next;
            q->front->next = p->next;
            (s->G[s->top]).num = p->data;
            (s->G[s->top]).time = d.time;
            delete p;
            q->count--;
            if (q->front->next == NULL) {
                q->rear = q-> front;
                return 1;
            }
        }
    }
    else {
        for(i = 0; i < (s->top); i++) {
            if (s->G[i].num != d.num)
                continue;
            else
                break;
        }
        if (i >= (s->top)) {
            cout<<"ERROR!\n";
            return -1;
        }
        x = d.time - (s->G[i]).time;
        y = fee * x;
        cout<<"停车时间为："<<x<<"小时，停车费用为："<<y<<"元!\n";
        k = new SqStack;
        k->s_top = -1;
        for (j= (s->top); j > i; j--) {
            k->s_top++;
            (k->H[k->s_top]).num = (s->G[j]).num;
            (k->H[k->s_top]).time = (s->G[j]).time;
            s->top--;
        }
        cout<<"临时栈中的信息为：（车号和时间）：\n";
        for (l = 0; l <= (k->s_top); l++) {
            cout<<(k->H[l]).num<<","<<(k->H[l]).time<<endl;
        }
        s->top--;
        while (k->s_top >= 0) {
            s->top++;
            (s->G[s->top]).state = 'A';
            (s->G[s->top]).num = (k->H[k->s_top]).num;
            (s->G[s->top]).time = (k->H[k->s_top]).time;
            k->s_top--;
        }
        if (q->count == 0) {
            cout<<"便道为空"<<endl;
            return  2;
        } else {
            s->top++;
            p = q->front->next;
            q->front->next = p->next;
            (s->G[s->top]).num = p->data;
            (s->G[s->top]).time = d.time;
            delete p;
            q->count--;
            if (q->front->next == NULL)
                q->rear = q-> front;
            return 3;

        }
    }
    return 0;
}

void Judge_Output (ParkStack *s, LinkQueue *q, car r) {
    if (r.state == 'E' || r.state == 'e')
        cout<<"STOP!"<<endl;
    else if (r.state == 'P' || r.state == 'p')
        cout<<"停车场中汽车数量为："<<(s->top) + 1<<endl;
    else if (r.state == 'W' || r.state == 'w')
        cout<<" 便道中汽车数量为："<<q->count<<endl;
    else if (r.state == 'A' || r.state == 'a')
        Append_cars(s,q,r);
    else if (r.state == 'D' || r.state == 'd')
        D_cars(s,q,r);
    else
        cout<<"ERROR!"<<endl;
}

int main() {
    ParkStack *s = new ParkStack;;
    LinkQueue *q = new LinkQueue;
    QNODE *p = new QNODE;
   car aa[MAXSIZE];
    int i;
    s->top = -1;
    p->next = NULL;
    q->front = q->rear = p;
    q->count = 0;
    cout<<"停 车 场 管 理 系 统"<<endl;
    cout<<"******************************************************************************\n";
    cout<<"*                                                                            *"<<endl;
    cout<<"*      A(a)车辆到达；D(d)车辆离开;P(p)停车场车辆总数;W(w)便道车辆总数;E(e)退出      *"<<endl;
    cout<<"*                                                                            *"<<endl;
    cout<<"******************************************************************************\n";
    for (i = 0; i < MAXSIZE; i++) {
        cout<<"请输入汽车的状态，车牌号和时间："<<endl;
        cin>>aa[i].state>>aa[i].num>>aa[i].time;
        Judge_Output(s,q,aa[i]);
        if (aa[i].state == 'E' || aa[i].state == 'e' ){
            break;
        }
    }
    return 0;
}