/*
String Permutations  Share on LinkedIn

Description:

Write a program to print out all the permutations of a string in alphabetical order.

Input sample:

The first argument will be a text file containing an input string, one per line. e.g. 

hat
Output sample:

Print to stdout, permutations of the string, comma separated, in alphabetical order.
e.g.

aht,ath,hat,hta,tah,tha
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
    vector<string> table;
    
    while (getline(ifs, line)) {
    	string cur = line;
    	table.clear();

    	table.push_back(cur);
    	next_permutation(cur.begin(), cur.end());
    	while (cur != line) {
    	 	table.push_back(cur);
    	 	next_permutation(cur.begin(), cur.end());
    	}

    	sort(table.begin(), table.end());

    	cout << table[0];
    	for (int i = 1; i < table.size(); ++i) {
    		cout << "," << table[i];
    	}
    	cout << endl;
    }

    return 0;
}