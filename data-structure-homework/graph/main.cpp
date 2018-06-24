//
// Created by yang on 08/06/2017.
//
#include <iostream>
#define MAX_VERTEX_NUM 64
using namespace std;

typedef struct edgenode {
    int adjvex;               // 该边所指向的顶点的序号
    float weight;             // 边上的权值
    struct edgenode * next;   // 指向下一条弧的指针
}*edgeptr;

typedef struct {
    int vertex;
    edgeptr link;
} vexnode;

typedef vexnode Adj_List[MAX_VERTEX_NUM];

int dvisited[MAX_VERTEX_NUM];

void dfs(Adj_List g, int v0) {
    edgeptr p;
    printf("%d ",v0);
    dvisited[v0]=1;
    p=g[v0].link;
    while( p != NULL) {
        if( dvisited[p->adjvex] == 0)
            dfs(g, p->adjvex);
        p = p->next;
    }
}

int bvisited[MAX_VERTEX_NUM];

void bfs(Adj_List g, int v0) {
    int f, r, v, q[MAX_VERTEX_NUM];
    edgeptr p = new edgenode;
    bvisited[v0]=1;
    printf("%d ", v0);
    f=0;
    r=0;
    p=g[v0].link;
    do {
        while (p != NULL) {
            v = p->adjvex;
            if (bvisited[v]==0) {
                q[r]=v;
                r++;
                printf("%d ", v);
                bvisited[v]=1;
            }
            p=p->next;  //找某一顶点的所有邻接点并进队
        }
        if (f!=r) { //v出队
            v=q[f];
            f++;
            p=g[v].link;
        }
    }while ((p!=NULL) || (f!=r));
}

int main() {
    Adj_List ga;

    int n, e;
    cout<<"请输入顶点数n,边数e"<<endl;
    cin>>n>>e;
    for(int i = 1; i <= n; i++) {
        ga[i].vertex = i;
        ga[i].link = NULL;
    }

    for (int k = 0; k < e ;k++) {
        int i, j;
        cout<<"请输入顶点对<i,j>"<<endl;
        cin>>i>>j;
        edgeptr p = new struct edgenode;
        p->adjvex = j;
        p->next = ga[i].link;
        ga[i].link = p;
    }
    dfs(ga, 1);
    cout<<endl;
    bfs(ga,1);

    return 0;
}