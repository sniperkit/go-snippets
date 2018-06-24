/*
Climbing Stairs  Share on LinkedIn

Description:

You are climbing a stair case. It takes n steps to reach to the top. Each time you can either climb 1 or 2 steps. In how many distinct ways can you climb to the top?

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file contains a positive integer which is the total number of stairs. e.g.

10
20
Output sample:

Print out the number of ways to climb to the top of the staircase. e.g.

89
10946
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    long cur;
    vector<long> table;
    table.push_back(0);
    table.push_back(1);
    table.push_back(2);

    while (ifs >> cur) {
    	while (table.size() <= cur) {
    		int len = table.size();
    		table.push_back(table[len-1]+table[len-2]);
    	}
    	cout << table[cur] << endl;
    }

    return 0;
}