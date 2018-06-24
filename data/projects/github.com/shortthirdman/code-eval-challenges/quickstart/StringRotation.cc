/*
String Rotation  Share on LinkedIn

Description:

You are given two strings. Determine if the second string is a rotation of the first string.

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file contains two comma separated strings. e.g.

Hello,lloHe
Basefont,tBasefon
Output sample:

Print out True/False if the second string is a rotation of the first. e.g.

True
True
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
    	string left;
    	string right;
    	replace(line.begin(), line.end(), ',', ' ');
        istringstream iss(line);
        iss >> left >> right;
        if (left.size() != right.size()) {
        	cout << "False" << endl;
        	continue;
        }

        bool found = false;
        for (int i = 0; i < left.size(); ++i) {
        	string temp(left);
        	rotate(temp.begin(), temp.begin()+i, temp.end());
        	if (temp.compare(right) == 0) {
        		cout << "True" << endl;
        		found = true;
        		break;
        	}
        }

        if (!found) cout << "False" << endl;
    }

    return 0;
}