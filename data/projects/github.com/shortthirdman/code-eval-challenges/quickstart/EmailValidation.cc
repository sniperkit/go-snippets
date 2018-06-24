/*
Email Validation  Share on LinkedIn

Description:

You are given several strings that may/may not be valid emails. You should write a regular expression that determines if the email id is a valid email id or not. You may assume all characters are from the english language.

Input sample:

Your program should accept as its first argument a filename. This file will contain several text strings, one per line. Ignore all empty lines. eg.

foo@bar.com
this is not an email id
admin#codeeval.com
good123@bad.com
Output sample:

Print out 'true' (all lowercase) if the string is a valid email. Else print out 'false' (all lowercase) .e.g.

true
false
false
true
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
        if (line.size() == 0) {
        	cout << "false" << endl;
        	continue;
        }

        bool at = false;
        bool dot = false;
        bool valid = true;
        for (int i = 0; i < line.size(); ++i) {
        	if (line[i] == '@') {
        		if (at || dot)	{
        			valid = false;
        			break;
        		} else {
        			at = true;
        		}
        	} else if (line[i] == '.') {
        		if (!at || dot || i == line.size()-1 || (i > 0 && line[i-1] == '@')) {
        			valid = false;
        			break;
        		} else {
        			dot = true;
        		}
        	} else {
        		if ( (line[i] >= 'a' && line[i] <= 'z') || (line[i] >= 'A' && line[i] <= 'Z') || (line[i] >= '0' && line[i] <= '9') ) {
        			continue;
        		} else {
        			valid = false;
        			break;
        		}
        	}
        }

        if (valid && at && dot) cout << "true" << endl;
        else cout << "false" << endl;
    }

    return 0;
}