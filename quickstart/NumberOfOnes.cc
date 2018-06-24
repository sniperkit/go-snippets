/*
Number of Ones  Share on LinkedIn

Description:

Write a program to determine the number of 1 bits in the internal representation of a given integer.

Input sample:

The first argument will be a text file containing an integer, one per line. e.g. 

10
22
56
Output sample:

Print to stdout, the number of ones in the binary form of each number.
e.g.

2
3
3
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
    
    while (getline(ifs, line)) {
        istringstream iss(line);
        iss >> N;
        int count = 0;
        while (N > 0) {
        	if (N&0x01 == 1) count++;
        	N >>= 1;
        }
        cout << count << endl;
    }

    return 0;
}