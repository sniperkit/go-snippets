/*
Number Pairs  Share on LinkedIn

Description:

You are given a sorted array of positive integers and a number 'X'. Print out all pairs of numbers whose sum is equal to X. Print out only unique pairs and the pairs should be in ascending order

Input sample:

Your program should accept as its first argument a filename. This file will contain a comma separated list of sorted numbers and then the sum 'X', separated by semicolon. Ignore all empty lines. If no pair exists, print the string NULL eg.

1,2,3,4,6;5
2,4,5,6,9,11,15;20
1,2,3,4;50
Output sample:

Print out the pairs of numbers that equal to the sum X. The pairs should themselves be printed in sorted order i.e the first number of each pair should be in ascending order .e.g.

1,4;2,3
5,15;9,11
NULL
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
    string line;
    vector<int> table;
    set<int> set;
    int target;
    
    while (getline(ifs, line)) {
    	int pos = line.find(';');
    	string remain = line.substr(pos+1);
    	istringstream(remain) >> target;

    	string data = line.substr(0, pos);
    	replace(data.begin(), data.end(), ',', ' ');
        istringstream iss(data);
        int temp = 0;
        table.clear();
        set.clear();
        while (iss >> temp) {
        	table.push_back(temp);
        	set.insert(temp);
        }

        sort(table.begin(), table.end());

        int found = 0;
        for (int i = 0; i < table.size(); ++i) {
        	if (set.find(target-table[i]) != set.end()) {
        		if (table[i] > target/2) break;
        		if (table[i] == target-table[i]) {
        			if (i+1 >= table.size() || table[i+1] != target-table[i]) break;
        		}

        		if (found == 0) found = 1;
        		else cout << ";";
        		cout << table[i] << "," << target-table[i];
        	}
        }
        if (found == 1) cout << endl;
        else cout << "NULL" << endl;
    }

    return 0;
}