/*
Multiples of a Number  Share on LinkedIn

Description:

Given numbers x and n, where n is a power of 2, print out the smallest multiple of n which is greater than or equal to x. Do not use division or modulo operator.

Input sample:

The first argument will be a text file containing a comma separated list of two integers, one list per line. e.g. 

13,8
17,16
Output sample:

Print to stdout, the smallest multiple of n which is greater than or equal to x, one per line.
e.g.

16
32
*/

#include <iostream>
#include <fstream>
#include <sstream>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    long x, n;
    string line;
    
    while (getline(ifs, line)) {
        for (int i = 0; i < line.size(); ++i)
            if (line[i] == ',') line[i] = ' ';

        istringstream iss(line);
        iss >> x >> n;

        while (n < x)
        	n <<= 1;
       	cout << n << endl;
    }

    return 0;
}