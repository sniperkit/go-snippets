/*
Sum of Integers from File  Share on LinkedIn

Description:

Print out the sum of integers read from a file.

Input sample:

The first argument to the program will be a text file containing a positive integer, one per line. e.g. 

5
12
NOTE: For solutions in JavaScript, assume that there are 7 lines of input
Output sample:

Print out the sum of all the integers read from the file. 
e.g.

17
*/

#include <iostream>
#include <fstream>
#include <sstream>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    int A;
    string line;
    
    long long sum = 0;
    while (getline(ifs, line)) {
        istringstream iss(line);
        iss >> A;
        sum += A;
    }

    cout << sum << endl;

    return 0;
}