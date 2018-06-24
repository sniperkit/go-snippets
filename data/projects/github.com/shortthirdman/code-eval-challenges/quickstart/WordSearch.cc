/*
Word Search  Share on LinkedIn

Description:

Given a 2D board and a word, find if the word exists in the grid. The word can be constructed from letters of sequentially adjacent cell, where adjacent cells are those horizontally or vertically neighboring. The same letter cell may not be used more than once.

Input sample:

The board to be used may be hard coded as:

[
[ABCE],
[SFCS],
[ADEE]
]
Your program should accept as its first argument a path to a filename. Each line in this file contains a word. e.g.
ASADB
ABCCED
Output sample:

Print out True if the word exists in the board, False otherwise. e.g.

False
True
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>

using namespace std;

bool find(vector<vector<bool> > mark, vector<string> &dict, int i, int j, string line, int pos) {
	if (pos >= line.size()) return true;
	if (i < 0 || i >= dict.size() || j < 0 || j >= dict[0].size()) return false;
	if (dict[i][j] != line[pos]) return false;
	if (mark[i][j] == true) return false;

	mark[i][j] = true;
	return find(mark, dict, i+1, j, line, pos+1) || 
		   find(mark, dict, i-1, j, line, pos+1) ||
		   find(mark, dict, i, j+1, line, pos+1) ||
		   find(mark, dict, i, j-1, line, pos+1);
}

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    vector<string> dict;
    dict.push_back("ABCE");
    dict.push_back("SFCS");
    dict.push_back("ADEE");
    
    while (getline(ifs, line)) {
        bool res = false;
        for (int i = 0; i < dict.size(); ++i) {
        	if (res == true) break;
        	int pos = dict[i].find(line[0]);
        	vector<vector<bool> > mark (3, vector<bool> (4, false));
        	while (pos != string::npos) {
        		res |= find(mark, dict, i, pos, line, 0);
				if (res == true) break;
        		pos = dict[i].find(line[0], pos+1);
        	}
        }

        if (res) cout << "True" << endl;
        else cout << "False" << endl; 
    }

    return 0;
}