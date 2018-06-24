/*
Reverse Groups  Share on LinkedIn

Description:

Given a list of numbers and a positive integer k, reverse the elements of the list, k items at a time. If the number of elements is not a multiple of k, then the remaining items in the end should be left as is.

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file contains a list of numbers and the number k, separated by a semicolon. The list of numbers are comma delimited. e.g.

1,2,3,4,5;2
1,2,3,4,5;3
Output sample:

Print out the new comma separated list of numbers obtained after reversing. e.g.

2,1,4,3,5
3,2,1,4,5
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    int k;
    string line;
    vector<int> table;

    while (getline(ifs, line)) {
    	int pos = line.find(';');
    	istringstream(line.substr(pos+1)) >> k;
    	string data = line.substr(0, pos);
    	replace(data.begin(), data.end(), ',', ' ');
        istringstream iss(data);
        int cur;
        table.clear();
        bool start = false;
        while (iss >> cur) {
        	table.push_back(cur);
        	if (table.size() == k) {
        		for (int i = table.size()-1; i >= 0; --i) {
        			if (start) cout << ",";
        			cout << table[i];
        			start = true;
        		}
        		table.clear();
        	}
        }
        for (int i = 0; i < table.size(); ++i) {
        	if (start) cout << ',';
        	cout << table[i];
        	start = true;
        }
        cout << endl;
    }

    return 0;
}