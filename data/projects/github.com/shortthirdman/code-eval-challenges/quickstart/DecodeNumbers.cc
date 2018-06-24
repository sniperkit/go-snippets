/*
Decode Numbers  Share on LinkedIn

Description:

You are given an encoded message containing only numbers. You are also provided with the following mapping:

A : 1
B : 2
C : 3
...
Z : 26
Given an encoded message, count the number of ways it can be decoded.
Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file is a testcase and contains an encoded message of numbers. e.g.

12
123
You may assume that the test cases contain only numbers.
Output sample:

Print out the different number of ways it can be decoded. e.g.

2
3
NOTE: 12 could be decoded as AB(1 2) or L(12). Hence the number of ways to decode 12 is 2.

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
    
    while (getline(ifs, line)) {
    	if (line.size() == 0) continue;

    	int table[line.size()+1];
    	table[line.size()] = 1;

    	for (int i = line.size()-1; i >= 0; --i) {
    		if (line[i] == '0') table[i] = 0;
    		else table[i] = table[i+1];

    		if ( i != line.size()-1 && 
    			 (line[i] == '1' || (line[i] == '2' && line[i+1] < '7')))
    			table[i] += table[i+2];
    	}

        cout << table[0] << endl;
    }

    return 0;
}