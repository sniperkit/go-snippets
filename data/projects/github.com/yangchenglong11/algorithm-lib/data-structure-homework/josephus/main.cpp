//
// Created by yang on 03/06/2017.
//
#include <iostream>

using namespace std;

void Josephus (int A [], int n, int s, int m) {
    int i, j, k, tmp;
    if (m == 0) {
        cout<<"m = 0 是无效的参数！\n";
        return;
    }

    for (i = 0; i < n; i++)
        A[i] = i + 1;

    i = s - 1;
    for (k = n; k > 1; k--) {
        if (i == k)
            i = 0;
        i = (i + m - 1) % k;
        if (i != k - 1) {
            tmp = A[i];
            for (j = i; j < k - 1; j++)
                A[j] = A [j + 1];
            A[k - 1] = tmp;
        }
    }

    for (k = 0; k < n/2; k++) {
        tmp = A[k];
        A[k] = A[n-k-1];
        A[n-k-1] = tmp;
    }

    for (i = 0; i < n; i++)
        cout<<A[i];
}

struct ListNode {
    int data;
    struct ListNode *next;
};

ListNode *createList (int n) {
    int i;
    ListNode *head = new ListNode;
    ListNode *r = new ListNode;

    for (i = 0; i < n; i++) {
        ListNode *s = new ListNode;
        s->data = i + 1;
        s->next = NULL;
        if ( i == 0 ) {
            head = s;
            r = s;
        }
        else {
            r->next = s;
            r = s;
            if (i == n-1)
                r->next = head;
        }
    }

    return head;
}

void display ( ListNode * head, int n, int m) {
    for (int k = 0; k < n; k++) {
        ListNode *p;
        if (m > 1) {
            p = head->next;
            for (int i = 1; i < m - 1; i++) {
                head = head->next;
                p = p -> next;
            }

            cout<<p->data;
            head->next = p->next;
            head = head->next;
            delete(p);
        }
        else if (m == 1) {
            p = head;
            cout<<p->data;
            head = p->next;
            delete(p);
        }
    }
}

void print(ListNode *pHead, int i)
{
    if (!pHead)
        return;
    ListNode *pre2, *cur2;
    pre2 = pHead;
    cur2 = pHead->next;
    if (i == 1)
    {
        pHead = pHead->next;
    }
    else
    {
        cout << pHead->data;
    }
    while (cur2 != pHead)
    {
        if (cur2->data == i)
        {
            pre2->next = cur2->next;
            free(cur2);
            cur2 = pre2->next;
        }
        cout << cur2->data;
        pre2->next = cur2->next;
        cur2 = pre2->next;
    }
}

int main() {
    int a[8];
    Josephus(a, 8, 1, 4);

    cout<<endl;

    ListNode * b = createList(8);
    display(b, 8, 4);

    return 0;
}