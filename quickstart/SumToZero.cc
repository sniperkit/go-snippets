/*
Sum to Zero  Share on LinkedIn

Description:

You are given an array of integers. Count the numbers of ways in which the sum of 4 elements in this array results in zero

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file consist of comma separated positive and negative integers. e.g.

2,3,1,0,-4,-1
0,-1,3,-2
Output sample:

Print out the count of the different number of ways that 4 elements sum to zero. e.g.

2
1
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    vector<int> table;
    string line;
    
    while (getline(ifs, line)) {
    	replace(line.begin(), line.end(), ',', ' ');
        istringstream iss(line);
        int cur;
        table.clear();
        while (iss >> cur) table.push_back(cur);
        sort(table.begin(), table.end());

        int count = 0;
        for (int i = 0; i < table.size()-2; ++i) {
        	for (int j = i+1; j < table.size()-1; ++j) {
        		int left = j+1;
        		int right = table.size()-1;

        		while (left < right) {
        			int val = table[i] + table[j] + table[left] + table[right];
        			if (val == 0) {
        				count++;
        				left++;
        				right--;
        			} else if (val < 0) {
        				left++;
        			} else {
        				right--;
        			}
        		}
        	}
        }

        cout << count << endl;

    }

    return 0;
}