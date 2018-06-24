/*
First Non-Repeated Character  Share on LinkedIn

Description:

Write a program to find the first non repeated character in a string.

Input sample:

The first argument will be a text file containing strings. e.g. 

yellow
tooth
Output sample:

Print to stdout, the first non repeating character, one per line.
e.g.

y
h
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
    vector<int> table;

    while (getline(ifs, line)) {
    	if (line.size() == 0) continue;
        transform(line.begin(), line.end(), line.begin(), ::tolower);
        table.clear();
        table.resize(26);
        for (int i = 0; i < line.size(); ++i) {
        	if (line[i] >= 'a' && line[i] <= 'z') table[line[i]-'a']++;
        }

        for (int i = 0; i < line.size(); ++i) {
        	if (line[i] >= 'a' && line[i] <= 'z' && table[line[i]-'a'] == 1) {
        		cout << line[i] << endl;
        		break;
        	}
        }
    }

    return 0;
}