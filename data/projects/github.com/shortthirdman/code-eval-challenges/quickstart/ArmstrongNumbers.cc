/*
Armstrong Numbers  Share on LinkedIn

Description:

An Armstrong number is an n-digit number that is equal to the sum of the n'th powers of its digits. Determine if the input numbers are Armstrong numbers.

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file has a positive integer. e.g.

6
153
351
Output sample:

Print out True/False if the number is an Armstrong number or not e.g.

True
True
False
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <cmath>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    int N;
    string line;
    
    while (getline(ifs, line)) {
        istringstream iss(line);
        iss >> N;

        int temp = N;
        int count = 0;
        while (temp > 0) {
        	++count;
        	temp /= 10;
        }

        int sum = 0;
        temp = N;
        while (temp > 0) {
        	sum += pow((double)(temp%10), count);
        	temp /= 10;
        }

        if (sum == N) cout << "True\n";
        else cout << "False\n";

    }

    return 0;
}