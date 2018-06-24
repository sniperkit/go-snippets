/*
Happy Numbers  Share on LinkedIn

Description:

A happy number is defined by the following process. Starting with any positive integer, replace the number by the sum of the squares of its digits, and repeat the process until the number equals 1 (where it will stay), or it loops endlessly in a cycle which does not include 1. Those numbers for which this process ends in 1 are happy numbers, while those that do not end in 1 are unhappy numbers.

Input sample:

The first argument is the pathname to a file which contains test data, one test case per line. Each line contains a positive integer. Each line is in the format: N i.e. a positive integer eg.

1
7
22
Output sample:

If the number is a happy number, print out a 1. If not, print out a 0 eg.

1
1
0
For the curious, here's why 7 is a happy number: 7->49->97->130->10->1. Here's why 22 is NOT a happy number: 22->8->64->52->29->85->89->145->42->20->4->16->37->58->89 ...
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <set>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    int N;
    string line;
    set<int> table;
    
    while (getline(ifs, line)) {
        istringstream iss(line);
        iss >> N;

        table.clear();
        while (1) {
        	int cur = 0;
        	while (N > 0) {
        		cur += (N%10)*(N%10);
        		N /= 10;
        	}
        	N = cur;
        	
        	if (table.find(N) != table.end()) {
        		cout << "0" << endl;
        		break;
        	}
        	if (N == 1) {
        		cout << "1" << endl;
        		break;
        	}
        	table.insert(N);
        }

    }

    return 0;
}