/*
Sum of Digits  Share on LinkedIn

Description:

Given a positive integer, find the sum of its constituent digits.

Input sample:

The first argument will be a text file containing positive integers, one per line. e.g. 

23
496
Output sample:

Print to stdout, the sum of the numbers that make up the integer, one per line.
e.g.

5
19
*/

#include <iostream>
#include <fstream>
#include <sstream>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    int A;
    string line;
    
    while (getline(ifs, line)) {
        istringstream iss(line);
        iss >> A;
        int sum = 0;
        while (A > 0) {
        	sum += A%10;
        	A /= 10;
        }
        cout << sum << endl;
    }

    return 0;
}