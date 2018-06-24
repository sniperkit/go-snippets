/*
Sum of integers  Share on LinkedIn

Description:

Write a program to determine the largest sum of contiguous integers in a list.

Input sample:

The first argument will be a text file containing a comma separated list of integers, one per line. e.g. 

-10, 2, 3, -2, 0, 5, -15
2,3,-2,-1,10
Output sample:

Print to stdout, the largest sum. In other words, of all the possible contiguous subarrays for a given array, find the one with the largest sum, and print that sum.
e.g.

8
12
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
        if (line.size() == 0) continue;
        replace(line.begin(), line.end(), ',', ' ');
        istringstream iss(line);
        table.clear();
        int temp;
        while (iss >> temp) table.push_back(temp);
        
        int sum = 0;
        int res = 0;
        for (int i = 0; i < table.size(); ++i) {
        	sum += table[i];
        	if (sum < 0) sum = 0;
        	res = max(res, sum);
        }

        cout << res << endl;
    }

    return 0;
}