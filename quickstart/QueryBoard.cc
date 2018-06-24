/*
Query Board  Share on LinkedIn

Description:

There is a board (matrix). Every cell of the board contains one integer, which is 0 initially.

The next operations can be applied to the Query Board:
SetRow i x: it means that all values in the cells on row "i" have been changed to value "x" after this operation.
SetCol j x: it means that all values in the cells on column "j" have been changed to value "x" after this operation.
QueryRow i: it means that you should output the sum of values on row "i".
QueryCol j: it means that you should output the sum of values on column "j".

The board's dimensions are 256x256
"i" and "j" are integers from 0 to 255
"x" is an integer from 0 to 31
Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file contains an operation of a query. E.g.

SetCol 32 20
SetRow 15 7
SetRow 16 31
QueryCol 32
SetCol 2 14
QueryRow 10
SetCol 14 0
QueryRow 15
SetRow 10 1
QueryCol 2
Output sample:

For each query, output the answer of the query. E.g.

5118
34
1792
3571
*/

#include <iostream>
#include <fstream>
#include <sstream>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    int board[256][256];

    for (int i = 0; i < 256; ++i)
    	for (int j = 0; j < 256; ++j)
    		board[i][j] = 0;
    
    while (getline(ifs, line)) {
        istringstream iss(line);
        string cmd;
        int pos;
        int val = 0;
        iss >> cmd >> pos;
        if (cmd == "SetCol") {
        	iss >> val;
        	for (int i = 0; i < 256; ++i)
        		board[i][pos] = val;
        } else if (cmd == "SetRow") {
        	iss >> val;
        	for (int i = 0; i < 256; ++i)
        		board[pos][i] = val;
        } else if (cmd == "QueryCol") {
        	for (int i = 0; i < 256; ++i)
        		val += board[i][pos];
        	cout << val << endl;
        } else if (cmd == "QueryRow") {
        	for (int i = 0; i < 256; ++i)
        		val += board[pos][i];
        	cout << val << endl;
        } else {
        	cout << "Wrong!" << endl;
        }
    }

    return 0;
}