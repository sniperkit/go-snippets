/*
Repeated Substring  Share on LinkedIn

Description:

You are to find the longest repeated substring in a given text. Repeated substrings may not overlap. If more than one substring is repeated with the same length, print the first one you find.(starting from the beginning of the text). NOTE: The substrings can't be all spaces

Input sample:

Your program should accept as its first argument a path to a filename.The input file contains several lines. Each line is one test case. Each line contains a test string. eg.

banana
Output sample:

For each set of input produce a single line of output which is the longest repeated substring. If there is none, print out the string NONE. eg.

an
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
    	bool found = false;

    	for (int len = line.size()/2; len > 0; --len) {
    		found = false;

        	for (int i = 0; i <= line.size()-len*2; ++i) {
        		string target = line.substr(i, len);

                int j = 0;
                while (target[j] == ' ') j++;
                if (j == target.size()) continue;

        		for (int k = i+len; k <= line.size()-len; ++k) {
        			string second = line.substr(k, len);
        			if (target.compare(second) == 0) {
        				cout << target << endl;
        				found = true;
        				break;
        			}
        		}
        		if (found == true) break;
    		}
            if (found == true) break;
        }

        if (found == false)
        	cout << "NONE" << endl;
    }

    return 0;
}