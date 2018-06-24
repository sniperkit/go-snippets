/*
Decimal To Binary  Share on LinkedIn

Description:

Given a decimal (base 10) number, print out its binary representation.

Input sample:

File containing positive whole decimal numbers, one per line. e.g. 

2
10
67
Output sample:

Print the decimal representation, one per line.
e.g.

10
1010
1000011
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
    string line;
    vector<int> table;

    while (getline(ifs, line)) {
        istringstream iss(line);
        iss >> N;
        table.clear();
        while (N > 0) {
        	table.push_back(N&0x01);
        	N >>= 1;
        }
        for (int i = table.size()-1; i >= 0; --i) 
        	cout << table[i];
        cout << endl;
    }

    return 0;
}