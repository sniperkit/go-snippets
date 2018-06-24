/*
String List  Share on LinkedIn

Description:

Credits: Challenge contributed by Max Demian.

You are given a number N and a string S. Print all of the possible ways to write a string of length N from the characters in string S, comma delimited in alphabetical order.

Input sample:

The first argument will be the path to the input filename containing the test data. Each line in this file is a separate test case. Each line is in the format: N,S i.e. a positive integer, followed by a string (comma separated) eg.

1,aa
2,ab
3,pop
Output sample:

Print all of the possible ways to write a string of length N from the characters in string S comma delimited in alphabetical order, with no duplicates. eg.

a
aa,ab,ba,bb
ooo,oop,opo,opp,poo,pop,ppo,ppp
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <set>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    
    while (getline(ifs, line)) {
        int pos = line.find(',');
        int len = atoi(line.substr(0, pos).c_str());

        line = line.substr(pos+1);
        set<char> dict;
        for (int i = 0; i < line.size(); ++i) {
        	dict.insert(line[i]);
        }

        vector<string> cur;
        vector<string> nxt;
        cur.push_back("");

        for (int i = 0; i < len; ++i) {
        	nxt.clear();
        	for (int j = 0; j < cur.size(); ++j) {
        		for (set<char>::iterator it = dict.begin(); it != dict.end(); ++it) {
        			nxt.push_back(cur[j]+*it);
        		}
        	}
        	swap(cur, nxt);
        }

        for (int i = 0; i < cur.size(); ++i) {
        	if (i != 0) cout << ",";
        	cout << cur[i];
        }
        cout << endl;
    }

    return 0;
}