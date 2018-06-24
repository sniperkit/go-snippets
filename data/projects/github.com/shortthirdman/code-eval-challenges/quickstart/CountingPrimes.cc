/*
Counting Primes  Share on LinkedIn

Description:

Given two integers N and M, count the number of prime numbers between N and M (both inclusive)

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file contains two comma separated positive integers. e.g.

2,10
20,30
Output sample:

Print out the number of primes between N and M (both inclusive)
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
    int N, M;
    vector<int> table;
    table.push_back(2);

    while (getline(ifs, line)) {
    	replace(line.begin(), line.end(), ',', ' ');
        istringstream iss(line);
        iss >> N >> M;

        for (int i = table.back()+1; i <= M; ++i) {
        	bool prime = true;
        	for (int j = 0; j < table.size() && table[j]*table[j] <= i; j++) {
        		if (i % table[j] == 0) {
        			prime = false;
        			break;
        		}
        	}
        	if (prime) table.push_back(i);
        }

        int count = 0;
        for (int i = 0; i < table.size(); i++) {
        	if (table[i] > M) break;
        	else if (table[i] < N) continue;
        	else ++count;
        }

        cout << count << endl;
    }

    return 0;
}