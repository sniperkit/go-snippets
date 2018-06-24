/*
Unique Elements  Share on LinkedIn

Description:

You are given a sorted list of numbers with duplicates. Print out the sorted list with duplicates removed.

Input sample:

File containing a list of sorted integers, comma delimited, one per line. e.g. 

1,1,1,2,2,3,3,4,4
2,3,4,5,5
Output sample:

Print out the sorted list with duplicates removed, one per line
e.g.

1,2,3,4
2,3,4,5
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <algorithm>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    
    while (getline(ifs, line)) {
        replace(line.begin(), line.end(), ',', ' ');
        istringstream iss(line);
        int cur = 0;
        iss >> cur;
        cout << cur;
        int pre = cur;
        while (iss >> cur) {
        	if (pre == cur) continue;
        	cout << "," << cur;
        	pre = cur;
        }
        cout << endl;
    }

    return 0;
}