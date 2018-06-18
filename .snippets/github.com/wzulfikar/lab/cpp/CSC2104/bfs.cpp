// Thu, 27 Apr 2017 at 17:49:43 MYT

#include <iostream>
#include <vector>
#include <queue>

using namespace std;

#define INF (1<<30)
#define pii pair<int, int>

#include <cstdlib>

void pause () {
	// pause the loop to observe changes in graph
	cout << "Press enter to continue ..."; 
    cin.get(); 
}

void clear_screen()
{
#ifdef WINDOWS
    std::system("cls");
#else
    // Assume POSIX
    std::system ("clear");
#endif
}

vector< vector<char> > graph; // char graph[size][size]

// 				   U  R  D  L
int dr [] = 	 {-1, 0, 1, 0}; // vector<int> dr = {-1, 0, 1, 0};
int dc [] = 	 { 0, 1, 0, -1}; // vector<int> dc = { 0, 1, 0, -1};

vector< vector <int> > dist;

void generateGraph(int row, int col, char symbol) {
	graph.resize(row);
	dist.resize(row);
	for (int i=0; i < row; i++) {
		graph[i].assign(col, symbol);
		dist[i].assign(col, INF);
	}
}

template<class T>
void displayGraph(vector< vector<T> >grid_2d) {
	for (int i=0; i < grid_2d.size(); i++) {
		for (int j=0; j < grid_2d[0].size(); j++) {
			if (grid_2d[i][j] == INF) {
				cout << "# ";
			} else {
				cout << grid_2d[i][j] << " ";
			}
		}
		cout << endl;
	}
	cout << endl;
}

bool isInside(int nr, int nc) {
	if (nr >= 0 && nr < graph.size() && nc >= 0 && nc < graph[0].size()) {
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
void bfs(int sr, int sc) {
	dist[sr][sc] = 0;

	queue< pii > q;
	q.push(make_pair(sr, sc));

	while(!q.empty()){
		pii cur = q.front();
		q.pop();

		for (int i=0; i< sizeof(dr); i++) {
			int nr = cur.first + dr[i];
			int nc = cur.second + dc[i];

			if (isInside(nr, nc)) {
				int new_distance = dist[cur.first][cur.second] + 1;
				if (new_distance < dist[nr][nc]) {
					dist[nr][nc] = new_distance;
					q.push(make_pair(nr, nc));

					pause();
				    clear_screen();

					displayGraph<int>(dist);
				}
			}
		}
	}
	// if (graph[sr][sc] != '#') return; // base case

	// graph[sr][sc] = depth;

	// for (int i = 0; i < sizeof(dr); i++) {
	// 	int nr = sr + dr[i];
	// 	int nc = sc + dc[i];

	// 	if (isInside(nr, nc)) {
	// 		dfs(nr, nc, depth + 1);
	// 	}
	// }
}

int main (){
	generateGraph(5, 5, '#');
	displayGraph<char>(graph);
	displayGraph<int>(dist);

	bfs(2, 2);

	displayGraph<int>(dist);
	return 0;
}