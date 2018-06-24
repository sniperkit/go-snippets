/*
Trailing String  Share on LinkedIn

Description:

You are given two strings 'A' and 'B'. Print out a 1 if string 'B' occurs at the end of string 'A'. Else a zero.

Input sample:

The first argument is a file, containing two strings, comma delimited, one per line. Ignore all empty lines in the input file.e.g. 

Hello World,World
Hello CodeEval,CodeEval
San Francisco,San Jose
Output sample:

Print out 1 if the second string occurs at the end of the first string. Else zero. Do NOT print out empty lines between your output. 
e.g.

1
1
0
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
    
    while (getline(ifs, line)) {
    	if (line.size() == 0) continue;
        int pos = line.find(',');
        string left = line.substr(0, pos);
        string right = line.substr(pos+1);
        if (right.size() > left.size()) {
        	cout << "0" << endl;
        	continue;
        }

        int should_start = left.size()-right.size();
        string cut = left.substr(should_start);
        if (cut.compare(right) != 0) {
        	cout << "0" << endl;
        } else {
        	cout << "1" << endl;
        }
    }

    return 0;
}