/*
Pascals Triangle  Share on LinkedIn

Description:

A pascals triangle row is contructed by looking at the previous row and adding the numbers to its left and right to arrive at the new value. If either the number to its left/right is not present, substitute a zero in it's place. More details can be found here: http://en.wikipedia.org/wiki/Pascal's_triangle. e.g. A Pascals triangle upto a depth of 6 can be shown as:

                1
              1   1
            1   2   1
           1  3   3   1
         1  4   6   4   1
        1  5  10  10  5   1
Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file contains a positive integer which indicates the depth of the triangle (1 based). e.g.

6
Output sample:

Print out the resulting pascal triangle upto the requested depth in row major form e.g.

1 1 1 1 2 1 1 3 3 1 1 4 6 4 1 1 5 10 10 5 1
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    int N;
    vector<int> table[2];
    int cur = 0;
    int nxt = 1;

    while (ifs >> N) {
    	table[cur].clear();
    	table[nxt].clear();
    	cur = 0;
    	nxt = 1;
    	table[cur].push_back(1);
    	for (int i = 0; i < N; ++i) {
	   		for (int j = 0; j < table[cur].size(); ++j) {
	    		if (j == 0) table[nxt].push_back(table[cur][j]);
	    		else table[nxt].push_back(table[cur][j]+table[cur][j-1]);
	    		cout << table[cur][j] << " ";
	    	}
	    	table[nxt].push_back(1);
	    	swap(cur, nxt);
	    	table[nxt].clear();
	   	}
	   	cout << endl;
    }

    return 0;
}