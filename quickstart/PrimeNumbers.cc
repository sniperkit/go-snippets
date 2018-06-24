/*
Prime Numbers  Share on LinkedIn

Description:

Print out the prime numbers less than a given number N. For bonus points your solution should run in N*(log(N)) time or better. You may assume that N is always a positive integer.

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file is one test case. Each test case will contain an integer n < 4,294,967,295. eg.

10
20
100
Output sample:

For each line of input, print out the prime numbers less than N, in ascending order, comma delimited. (There should not be any spaces between the comma and numbers) eg.

2,3,5,7
2,3,5,7,11,13,17,19
2,3,5,7,11,13,17,19,23,29,31,37,41,43,47,53,59,61,67,71,73,79,83,89,97
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    long val;
    vector<long> table;
    table.push_back(2);

    while (ifs >> val) {
    	for (long i = table.back()+1; i < val; ++i) {
    		bool prime = true;
    		for (int j = 0; j < table.size() && table[j]*table[j] <= i; ++j) {
    			if (i % table[j] == 0) {
    				prime = false;
    				break;
    			}
    		}
    		if (prime) 
    			table.push_back(i);
    	}
    	cout << table[0];
    	for (int i = 1; i < table.size(); i++) {
    		if (table[i] >= val) break;
    		cout << ',' << table[i];
    	}
    	cout << endl;
    }

    return 0;
}