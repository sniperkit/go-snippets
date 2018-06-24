//
// Created by yang on 09/06/2017.
//
#include "iostream"

using namespace std;

typedef struct BiTNode {
    char data;
    struct BiTNode *lchild,*rchild;
} BiTNode,*BiTree;//定义结点类型

BiTree CreateBiTree() { //创建树
    char p;
    BiTree T;
    scanf("%c",&p);
    if(p == '0')
        T = NULL;
    else {
        T = (BiTNode *)malloc(sizeof(BiTNode));  //为结点开辟空间
        T->data=p;
        T->lchild = CreateBiTree();
        T->rchild = CreateBiTree();
    }

    return (T);
}
void PreOrder(BiTree T) {  //先序
    if(T!=NULL) {
        printf("%c",T->data);
        PreOrder(T->lchild);
        PreOrder(T->rchild);
    }
}
void InOrder(BiTree T) {   //中序
    if(T!=NULL)
    {
        InOrder(T->lchild);
        printf("%c",T->data);
        InOrder(T->rchild);
    }
}
void PostOrder(BiTree T) {  //后序
    if(T!=NULL) {
        PostOrder(T->lchild);
        PostOrder(T->rchild);
        printf("%c",T->data);
    }
}
int main() {
    BiTree Ta;
    Ta=CreateBiTree();

    printf("先序遍历:");
    printf("\n");
    PreOrder(Ta);
    printf("\n");

    printf("中序遍历:");
    printf("\n");
    InOrder(Ta);
    printf("\n");

    printf("后序遍历:");
    printf("\n");
    PostOrder(Ta);

    return 0;
}