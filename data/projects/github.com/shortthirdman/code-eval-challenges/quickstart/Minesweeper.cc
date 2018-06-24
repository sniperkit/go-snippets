/*
Minesweeper  Share on LinkedIn

Description:

You will be given an M*N matrix. Each item in this matrix is either a '*' or a '.'. A '*' indicates a mine whereas a '.' does not. The objective of the challenge is to output a M*N matrix where each element contains a number (except the positions which actually contain a mine which will remain as '*') which indicates the number of mines adjacent to it. Notice that each position has at most 8 adjacent positions e.g. left, top left, top, top right, right, ...

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file contains M,N, a semicolon and the M*N matrix in row major form. e.g.

3,5;**.........*...
4,4;*........*......
Output sample:

Print out the new M*N matrix (in row major form) with each position(except the ones with the mines) indicating how many adjacent mines are there. e.g.

**100332001*100
*10022101*101110
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
    vector<vector<char> > data;

    while (getline(ifs, line)) {
    	replace(line.begin(), line.end(), ',', ' ');
    	replace(line.begin(), line.end(), ';', ' ');
        istringstream iss(line);

        int N, M;
        iss >> N >> M;
        data.clear();
        for (int i = 0; i < N; ++i) {
        	vector<char> x;
        	for (int j = 0; j < M; ++j) {
        		char cur;
        		iss >> cur;
        		x.push_back(cur);
        	}
        	data.push_back(x);
        }

        for (int i = 0; i < N; ++i) {
        	for (int j = 0; j < M; ++j) {
        		if (data[i][j] == '*') continue;

        		int res = 0;
        		if (i > 0 	&& data[i-1][j]=='*') res++;
        		if (j > 0 	&& data[i][j-1]=='*') res++;
        		if (i < N-1 && data[i+1][j]=='*') res++;
        		if (j < M-1 && data[i][j+1]=='*') res++;

        		if (i > 0 	&& j > 0 	&& data[i-1][j-1]=='*') res++;
        		if (i > 0 	&& j < M-1 	&& data[i-1][j+1]=='*') res++;
        		if (i < N-1 && j > 0 	&& data[i+1][j-1]=='*') res++;
        		if (i < N-1 && j < M-1 	&& data[i+1][j+1]=='*') res++;

        		data[i][j] = '0'+res;
        	}
        }

        for (int i = 0; i < N; ++i) {
        	for (int j = 0; j < M; ++j) {
        		cout << data[i][j];
        	}
        }
        cout << endl;


    }

    return 0;
}