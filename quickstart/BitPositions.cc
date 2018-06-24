/*
Bit Positions  Share on LinkedIn

Description:

Given a number n and two integers p1,p2 determine if the bits in position p1 and p2 are the same or not. Positions p1,p2 and 1 based.

Input sample:

The first argument will be a text file containing a comma separated list of 3 integers, one list per line. e.g. 

86,2,3
125,1,2
Output sample:

Print to stdout, 'true'(lowercase) if the bits are the same, else 'false'(lowercase).
e.g.

true
false
*/

#include <iostream>
#include <fstream>
#include <sstream>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    long long n;
    int a, b;
    string line;
    
    while (getline(ifs, line)) {
        for (int i = 0; i < line.size(); ++i)
            if (line[i] == ',') line[i] = ' ';

        istringstream iss(line);
        iss >> n >> a >> b;

        int a_mask = 1<<(a-1);
        int b_mask = 1<<(b-1);

        if ( ((n&a_mask) > 0 && (n&b_mask) > 0) ||
        	 ((n&a_mask) == 0 && (n&b_mask) == 0) )
        	cout << "true" << endl;
        else
        	cout << "false" << endl;
    }

    return 0;
}