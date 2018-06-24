/*
String Searching  Share on LinkedIn

Description:

You are given two strings. Determine if the second string is a substring of the first (Do NOT use any substr type library function). The second string may contain an asterisk(*) which should be treated as a regular expression i.e. matches zero or more characters. The asterisk can be escaped by a \ char in which case it should be interpreted as a regular '*' character. To summarize: the strings can contain alphabets, numbers, * and \ characters.

Input sample:

File containing two comma delimited strings per line. e.g. 

Hello,ell
This is good, is 
CodeEval,C*Eval
Old,Young
Output sample:

If the second string is indeed a substring of the first, print out a 'true'(lowercase), else print out a 'false'(lowercase), one per line.
e.g.

true
true
true
false
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
    	int pos = line.find(',');
    	string left = line.substr(0, pos);
    	string right = line.substr(pos+1);

    	pos = right.find("\\*");
    	while (pos != string::npos) {
    		right.replace(pos, 2, "@");
    		pos = right.find("\\*");
    	}

    	vector<string> table;
    	pos = right.find('*');
    	while (pos != string::npos) {
    		string temp = right.substr(0, pos);
    		table.push_back(temp);
    		right = right.substr(pos+1);
    		pos = right.find('*');
    	}
    	if (right.size() > 0) {
    		table.push_back(right);
    	}

		reverse(table.begin(), table.end());
    	for (int i = 0; i < table.size(); ++i) {
    		replace(table[i].begin(), table[i].end(), '@', '*');
    	}

    	bool res = true;
    	while (table.size() > 0) {
    		string cur = table.back();
    		pos = left.find(cur);
    		if (pos == string::npos) { res = false; break; }
    		left = left.substr(pos+1);
    		if (left.size() == 0) { res = false; break; }
    		table.pop_back();
    	}

    	if (res) cout << "true" << endl;
    	else cout << "false" << endl;
    }

    return 0;
}