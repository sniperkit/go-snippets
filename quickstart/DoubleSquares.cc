/*
Double Squares  Share on LinkedIn

Description:

Credits: This challenge appeared in the Facebook Hacker Cup 2011.

A double-square number is an integer X which can be expressed as the sum of two perfect squares. For example, 10 is a double-square because 10 = 3^2 + 1^2. Your task in this problem is, given X, determine the number of ways in which it can be written as the sum of two squares. For example, 10 can only be written as 3^2 + 1^2 (we don't count 1^2 + 3^2 as being different). On the other hand, 25 can be written as 5^2 + 0^2 or as 4^2 + 3^2. 
NOTE: Do NOT attempt a brute force approach. It will not work. The following constraints hold: 
0 <= X <= 2147483647
1 <= N <= 100

Input sample:

You should first read an integer N, the number of test cases. The next N lines will contain N values of X.

5
10
25
3
0
1
Output sample:

e.g.

1
2
0
1
1
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <set>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    int N;
    ifs >> N;
    set<long long> set;
    vector<long long> table;
    table.push_back(0);
    set.insert(0);

    for (int i = 0; i < N; ++i) {
    	long long x;
    	ifs >> x;

    	while (table.back() < x) {
    		long long cur = table.size()*table.size();
    		table.push_back(cur);
    		set.insert(cur);
    	}

    	int count = 0;
    	for (int j = 0; j < table.size(); j++) {
    		if (table[j] > x/2) break;
    		if (set.find(x-table[j]) != set.end()) {
    			count++;
    		}
    	}
    	cout << count << endl;
    }

    return 0;
}