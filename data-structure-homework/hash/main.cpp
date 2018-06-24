//
// Created by yang on 19/06/2017.
//
#include "iostream"
#include "fstream"
#include "stdio.h"
#define m 8

using namespace std;

typedef struct node {
    char name[12];
    char num[12];
    char address[12];
    struct node *next;
}*Link_list;

int hash_(char num[]) { //哈希函数
    int i = 0;
    int k = 0;
    while(num[i] != NULL) {
        k += num[i] - '0';
        i++;
    }
    k = k % 7;
    return(k);
}
void hashsearch(Link_list user[], char n[]) {  //查找
    int k1;
    k1 = hash_(n);
    Link_list p = user[k1];
    if(p->next == NULL)
        cout<<"没有此条记录"<<endl;
    else {
        p = p->next;
        while (strcmp(n, p->num) != 0 && p->next != NULL)
            p = p->next;
        if(strcmp(n, p->num) == 0)
            cout <<"姓名:"<<p->name<<"  电话号码:"<<p->num<<"  地址:"<<p->address<< endl;
        else
            cout<<"没有此条记录"<<endl;
    }
}

void hashinsert(Link_list user[], Link_list q) {    //插入
    int k = hash_(q->num);
    Link_list s = user[k];
    if(user[k]->next == NULL) {
        Link_list p = new node;
        strcpy(p->name, q->name);
        strcpy(p->num, q->num);
        strcpy(p->address, q->address);
        p->next = NULL;
        user[k]->next = p;
    }
    else {
        while(strcmp(s->num, q->num) != 0 && s->next != NULL)
            s = s->next;
        if(strcmp(s->num, q->num) == 0)
            cout<<"此记录已存在"<<endl;
        else {
            Link_list p = new node;
            strcpy(p->name, q->name);
            strcpy(p->num, q->num);
            strcpy(p->address, q->address);
            p->next = NULL;
            s->next = p;
        }
    }
}

void hashdelete(Link_list user[], char n[]) { //删除
    int k1;
    k1 = hash_(n);
    if(user[k1]->next == NULL)
        cout<<"该记录不存在"<<endl;
    else {
        Link_list s = user[k1];
        Link_list p;
        int r;  //标志,r==1时找到
        do {
            r = 1;
            p = s;
            s = s->next;
            if(strcmp(s->num, n) != 0)
                r = 0;
        }while((r == 0) && (s->next != NULL));
        if(r == 0)
            cout<<"此记录不存在"<<endl;
        else {
            p->next = s->next;
            delete s;
            cout<<"删除成功"<<endl;
        }
    }

}
void show(Link_list user[]) {  //显示哈希表
    int i;
    Link_list p;
    for (i = 0; i < 8; i++)
    {
        p = user[i]->next;
        while (p)
        {
            cout <<"姓名:"<<p->name<<"  电话号码:"<<p->num<<"  地址:"<<p->address<< endl;
            p = p->next;
        }
    }
}
void menu() {  //菜单
    cout<<"请输入要执行的功能："<<endl;
    cout<<" 1.插入记录" <<endl;
    cout<<" 2.查找记录" <<endl;
    cout<<" 3.删除记录" <<endl;
    cout<<" 4.显示记录" <<endl;
}

int main() {
    Link_list user[m];
    for(int i = 0; i <= m; i++)
    {
        user[i] = new node;
        user[i]->next = NULL;
    }
    FILE *fp;
    fp = fopen("/Users/yang/CLionProjects/hash/hash.txt","r");
    if(!fp)
        cout<<"File open error!\n";

    Link_list w = new node;
    for(int i = 0;i < 3;)
    {
        if(fscanf(fp,"%s%s%s", w->name, w->num, w->address) == 3)
        {
            hashinsert(user, w);
            i++;
        }
    }
    fclose(fp);
    cout<<"数据已录入"<<endl;
    int choose;
    char t[12];
    Link_list t1 = new node;
    t1->next = NULL;
    loop:
    menu();
    cin>>choose;
    switch(choose)
    {
        case 1:
            cout<<"请输入要插入的数据：";
            cin>>t1->name>>t1->num>>t1->address;
            hashinsert(user, t1);
            goto loop;
        case 2:
            cout<<"请输入要查找的手机号码：";
            cin>>t;
            hashsearch(user, t);
            goto loop;
        case 3:
            cout<<"请输入要删除的手机号码：";
            cin>>t;
            hashdelete(user, t);
            goto loop;
        case 4:
            show(user);
            goto loop;
    }

    return 0;
}