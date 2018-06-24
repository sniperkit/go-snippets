/*
Spiral Printing  Share on LinkedIn

Description:

Write a program to print a 2D array (n x m) in spiral order (clockwise)

Input sample:

Your program should accept as its first argument a path to a filename.The input file contains several lines. Each line is one test case. Each line contains three items (semicolon delimited). The first is 'n'(rows), the second is 'm'(columns) and the third is a single space separated list of characters/numbers in row major order. eg.

3;3;1 2 3 4 5 6 7 8 9
Output sample:

Print out the matrix in clockwise fashion, one per line, space delimited. eg.

1 2 3 6 9 8 7 4 5
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
    vector<vector<string> > table;

    while (getline(ifs, line)) {
    	replace(line.begin(), line.end(), ';', ' ');
        istringstream iss(line);

        int N, M;
        iss >> N >> M;

        table.clear();
        for (int i = 0; i < N; ++i) {
        	vector<string> x;
        	for (int j = 0; j < M; ++j) {
        		string cur;
        		iss >> cur;
        		x.push_back(cur);
        	}
        	table.push_back(x);
        }

        cout << table[0][0];
        int layer = 0;
        int direction = 0;
        int i = 0;
        int j = 1;
        int counter = 1;
        while (1) {
        	cout << " " << table[i][j];
        	if (++counter == N*M) { break; }
			switch (direction) {
                case 0 :
                    if (j+1 == M-layer) {
                        direction = 1;
                        i++;
                    } else { j++; }
                    break;
                case 1 :
                    if (i+1 == N-layer) {
                        direction = 2;
                        j--;
                    } else { i++; }
                    break;
                case 2 :
                    if (j == layer) {
                        direction = 3;
                        i--;
                    } else { j--; }
                    break;
                case 3 :
                    if (i-1 == layer) {
                        layer++;
                        direction = 0;
                        j++;
                    } else { i--; }
                    break;
                default :
                    break;
            }
        }

        cout << endl;
    }

    return 0;
}