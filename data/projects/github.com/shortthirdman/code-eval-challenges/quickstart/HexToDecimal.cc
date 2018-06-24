/*
Hex to Decimal  Share on LinkedIn

Description:

You will be given a hexadecimal(base 16) number. Convert it into decimal (base 10)

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file contains a hex number. You may assume that the hex number does not have the leading 'Ox'. Also all alpha characters(e.g. a through f) in the input will be in lowercase e.g.

9f
11
Output sample:

Print out the equivalent decimal number e.g.

159
17
*/

#include <iostream>
#include <fstream>
#include <sstream>

using namespace std;

int helper(char c, int pos) {
	int m = 1;
	while (pos-- > 0) m *= 16;

	int val = 0;
	if (c >= '0' && c <= '9') val = c - '0';
	else if (c >= 'a' && c <= 'f') val = c - 'a' + 10;

	return val * m;
}

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    
    while (getline(ifs, line)) {
    	int val = 0;
        for (int i = line.size()-1; i >= 0; --i) {
        	char c = line[i];
        	val += helper(c, line.size()-i-1);
        }
        cout << val << endl;
    }

    return 0;
}