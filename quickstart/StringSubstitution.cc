/*
String Substitution  Share on LinkedIn

Description:

Credits: This challenge was contributed by Sam McCoy

Given a string S, and a list of strings of positive length, F1,R1,F2,R2,...,FN,RN, proceed to find in order the occurrences (left-to-right) of Fi in S and replace them with Ri. All strings are over alphabet { 0, 1 }. Searching should consider only contiguous pieces of S that have not been subject to replacements on prior iterations. An iteration of the algorithm should not write over any previous replacement by the algorithm.

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file is one test case. Each test case will contain a string, then a semicolon and then a list of comma separated strings.eg.

10011011001;0110,1001,1001,0,10,11
Output sample:

For each line of input, print out the string after substitutions have been made.eg.

11100110
For the curious, here are the transitions for the above example: 10011011001 => 10100111001 [replacing 0110 with 1001] => 10100110 [replacing 1001 with 0] => 11100110 [replacing 10 with 11] => 11100110
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
    	int pos = line.find(";");
    	string target = line.substr(0, pos);

    	line = line.substr(pos+1);
    	replace(line.begin(), line.end(), ',', ' ');
        istringstream iss(line);
        string f, r;
        vector<string> table;
        vector<int> change(target.size(), 0);
        while (iss >> f >> r) {
        	table.push_back(f);
        	table.push_back(r);
        }

        int i = 0;
        while (i < table.size()) {
        	int pos = target.find(table[i]);
        	while (pos != string::npos) {
        		
        		int flag = 0;
        		for (int j = pos; j < pos+table[i].size(); ++j) {
        			flag += change[j];
        		}

        		if (flag == 0) {
        			target.replace(pos, table[i].size(), table[i+1]);
        			change.erase(change.begin()+pos, change.begin()+pos+table[i].size());
        			vector<int> x (table[i+1].size(), 1);
        			change.insert(change.begin()+pos, x.begin(), x.end());
        		}
        		pos = target.find(table[i], pos+1);
        	}
        	i += 2;
        }

        cout << target << endl;
    }

    return 0;
}