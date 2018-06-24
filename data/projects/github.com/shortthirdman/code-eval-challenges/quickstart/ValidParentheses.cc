/*
Valid parentheses  Share on LinkedIn

Description:

Given a string comprising just of the characters (,),{,},[,] determine if it is well-formed or not.

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file contains a string comprising of the characters mentioned above. e.g.

()
([)]
Output sample:

Print out True or False if the string is well-formed e.g.

True
False
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <stack>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    stack<char> s;
    
    while (getline(ifs, line)) {
        istringstream iss(line);
        char cur;
        while (!s.empty()) s.pop();
        bool valid = true;
        while (iss >> cur && valid) {
        	switch (cur) {
        	case '(':
        	case '[':
        	case '{':
        		s.push(cur);
        		break;
        	case ')':
				if (s.empty() || s.top() != '(') {
					valid = false;
					break;
				}
				s.pop();
				break;
			case ']':
				if (s.empty() || s.top() != '[') {
					valid = false;
					break;
				}
				s.pop();
				break;
			case '}':
				if (s.empty() || s.top() != '{') {
					valid = false;
					break;
				}
				s.pop();
				break;
			default:
				valid = false;
				break;
        	}
        }
        if (!s.empty() || !valid) cout << "False" << endl;
        else cout << "True" << endl;
    }

    return 0;
}