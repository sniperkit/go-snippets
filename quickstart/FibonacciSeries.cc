/*
Fibonacci Series  Share on LinkedIn

Description:

The Fibonacci series is defined as: F(0) = 0; F(1) = 1; F(n) = F(n-1) + F(n-2) when n>1;. Given a positive integer 'n', print out the F(n).

Input sample:

The first argument will be a text file containing a positive integer, one per line. e.g. 

5
12
Output sample:

Print to stdout, the fibonacci number, F(n).
e.g.

5
144
*/

#include <iostream>
#include <fstream>
#include <sstream>

using namespace std;

long fib(int A) {
	if (A < 3) return 1;

	long a = 1;
	long b = 1;
	long c = 2;

	while (A > 2) {
		c = a+b;
		a = b;
		b = c;
		A--;
	}

	return c;
}


int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    int A;
    string line;
    
    while (getline(ifs, line)) {
        istringstream iss(line);
        iss >> A;
        cout << fib(A) << endl;
    }

    return 0;
}