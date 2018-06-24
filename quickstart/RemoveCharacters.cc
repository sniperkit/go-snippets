/*
Remove Characters  Share on LinkedIn

Description:

Write a program to remove specific characters from a string.

Input sample:

The first argument will be a text file containing an input string followed by a comma and then the characters that need to be scrubbed. e.g. 

how are you, abc
hello world, def
Output sample:

Print to stdout, the scrubbed strings, one per line. Trim out any leading/trailing whitespaces if they occur. 
e.g.

how re you
hllo worl
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
    vector<char> table;

    while (getline(ifs, line)) {
        table.clear();
        int i = 0;
        while (line[i] != ',') {
        	table.push_back(line[i++]);
        }
        
        i += 2;  // skip , and space
        while (i < line.size()) {
        	char temp = line[i];
        	vector<char>::iterator it = find(table.begin(), table.end(), line[i]);
        	while (it != table.end()) {
        		table.erase(it);
        		it = find(table.begin(), table.end(), line[i]);
        	}
        	i++;
        }

        for (int i = 0; i < table.size(); ++i)
        	cout << table[i];
        cout << endl;
    }

    return 0;
}