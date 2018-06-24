/*
Distinct Subsequences  Share on LinkedIn

Description:

A subsequence of a given sequence S consists of S with zero or more elements deleted. Formally, a sequence Z = z1z2..zk is a subsequence of X = x1x2...xm, if there exists a strictly increasing sequence <i1,i2...ik> of indicies of X such that for all j=1,2,...k we have Xij = Zj. e.g. Z=bcdb is a subsequence of X=abcbdab with corresponding index sequence <2,3,5,7>

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file contains two comma separated strings. The first is the sequence X and the second is the subsequence Z. e.g.

babgbag,bag
rabbbit,rabbit
Output sample:

Print out the number of distinct occurrences of Z in X as a subsequence e.g.

5
3

*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>

using namespace std;

int helper(string a, string b, int i, int j) {
	if (j == b.size()) {
		if (i == a.size() && a[i-1] != b[j-1]) return 0;
		else return 1;
	} 
	if (i == a.size()) return 0;

	if (a[i] == b[j]) 
		return helper(a, b, i+1, j+1) + helper(a, b, i+1, j);
	else 
		return helper(a, b, i+1, j);
}

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    
    while (getline(ifs, line)) {
    	replace(line.begin(), line.end(), ',', ' ');
        istringstream iss(line);
        string a, b;
        iss >> a >> b;

        cout << helper(a, b, 0, 0) << endl;
    }

    return 0;
}