/*
N Mod M  Share on LinkedIn

Description:

Given two integers N and M, calculate N Mod M (without using any inbuilt modulus operator).

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file contains two comma separated positive integers. e.g.

20,6
2,3
You may assume M will never be zero.
Output sample:

Print out the value of N Mod M
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <algorithm>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    int N, M;
    string line;
    
    while (getline(ifs, line)) {
    	replace(line.begin(), line.end(), ',', ' ');
        istringstream iss(line);
        iss >> N >> M;

        while (N >= M) {
        	N -= M;
        }
        cout << N << endl;
    }

    return 0;
}