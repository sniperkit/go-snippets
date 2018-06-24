/*
Detecting Cycles  Share on LinkedIn

Description:

Given a sequence, write a program to detect cycles within it.
Input sample:

A file containing a sequence of numbers (space delimited). The file can have multiple such lines. e.g

2 0 6 3 1 6 3 1 6 3 1
Ensure to account for numbers that have more than one digit eg. 12. If there is no sequence, ignore that line.
Output sample:

Print to stdout the first sequence you find in each line. Ensure that there are no trailing empty spaces on each line you print. e.g.

6 3 1
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
        istringstream iss(line);
        if (line.size() == 0) continue;
        
        table.clear();
        int temp;
        while (iss >> temp)
            table.push_back(temp);

        int start = 0;
        int len = 0;
        for (int i = 0; i < table.size()-2; ++i) {
            for (int j = i+1; j < table.size()-1; ++j) {
                if (table[i] == table[j]) {
                    int k = 1;
                    while (i+k < j && j+k < table.size() && table[i+k] == table[j+k])
                        ++k;
                    if (k > len) {
                        len = k;
                        start = i;
                    }
                }
            }
        }

        cout << table[start];
        while (len-- > 1) {
            cout << " " << table[++start];
        }
        cout << endl;
    }

    return 0;
}