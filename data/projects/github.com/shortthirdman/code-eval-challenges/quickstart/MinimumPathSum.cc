/*
Minimum Path Sum  Share on LinkedIn

Description:

You are given an n*n matrix of integers. You can move only right and down. Calculate the minimal path sum from the top left to the bottom right

Input sample:

Your program should accept as its first argument a path to a filename. The first line will have the value of n(the size of the square matrix). This will be followed by n rows of the matrix. (Integers in these rows will be comma delimited). After the n rows, the pattern repeats. e.g.

2
4,6
2,8
3
1,2,3
4,5,6
7,8,9
Output sample:

Print out the minimum path sum for each matrix. e.g.

14
21
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    int N;
    vector<vector<int> > data;

    while (ifs >> N) {
    	data.clear();
    	for (int i = 0; i < N; ++i) {
    		getline(ifs, line);
    		if (line.size() == 0) {
    			--i;
    			continue;
    		}
    		replace(line.begin(), line.end(), ',', ' ');

    		int tmp;
    		vector<int> x;
    		istringstream iss(line);
    		while (iss >> tmp) {
    			x.push_back(tmp);
    		}
    		data.push_back(x);
    	}

		int table[N][N];
		table[0][0] = data[0][0];
		for (int i = 1; i < N; i++) {
			table[i][0] = table[i-1][0] + data[i][0];
			table[0][i] = table[0][i-1] + data[0][i];
		}

		for (int i = 1; i < N; ++i) {
			for (int j = 1; j < N; ++j) {
				table[i][j] = data[i][j] + min(table[i-1][j], table[i][j-1]);
			}
		}

		cout << table[N-1][N-1] << endl;
    }

    return 0;
}

class Solution {
public:
    int minPathSum(vector<vector<int> > &grid) {
        // Start typing your C/C++ solution below
        // DO NOT write int main() function
        int m = grid.size();
        if (m == 0) return 0;
        int n = grid[0].size();
        if (n == 0) return 0;
        
        int table[m][n];
        table[0][0] = grid[0][0];
        
        for (int i = 1; i < m; i++) table[i][0] = table[i-1][0] + grid[i][0];
        for (int i = 1; i < n; i++) table[0][i] = table[0][i-1] + grid[0][i];
        
        for (int i = 1; i < m; i++) {
            for (int j = 1; j < n; j++) {
                table[i][j] = grid[i][j] + min(table[i-1][j], table[i][j-1]);
            }
        }
        
        return table[m-1][n-1];
    }
};
