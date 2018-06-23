// Thu, 20 Apr 2017 at 17:42:03 MYT

#include <iostream>
#include <vector>

using namespace std;

vector< vector<char> > graph; // char grapch[size][size]

// 				   U  R  D  L
int dr [] = 	 {-1, 0, 1, 0}; // vector<int> dr = {-1, 0, 1, 0};
int dc [] = 	 { 0, 1, 0, -1}; // vector<int> dc = { 0, 1, 0, -1};

bool isInside(int r, int c) {
	if (r >= 0 && r < graph.size() && c >= 0 && c < graph[0].size()) {
		return true;
	}
	return false;
}

/**
 * Recursive function
 * 
 * @param sr [description]
 * @param sc [description]
 */
void dfs(int sr, int sc, char depth='A') {
	if (graph[sr][sc] != '#') return; // base case

	graph[sr][sc] = depth;

	for (int i = 0; i < sizeof(dr); i++) {
		int nr = sr + dr[i];
		int nc = sc + dc[i];

		if (isInside(nr, nc)) {
			dfs(nr, nc, depth + 1);
		}
	}
}

void generateGraph(int row, int col, char symbol) {
	graph.resize(row);
	for (int i=0; i < row; i++) {
		graph[i].assign(col, symbol);
	}
}

void displayGraph() {
	for (int i=0; i < graph.size(); i++) {
		for (int j=0; j < graph[i].size(); j++) {
			cout << graph[i][j] << " ";
		}
		cout << endl;
	}
	cout << endl;
}

int main (){
	generateGraph(5, 5, '#');
	displayGraph();

	dfs(4, 0);
	displayGraph();
}