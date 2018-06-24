/*
Lowercase  Share on LinkedIn

Description:

Given a string write a program to convert it into lowercase.

Input sample:

The first argument will be a text file containing sentences, one per line. You can assume all characters are from the english language. e.g. 

HELLO CODEEVAL
This is some text
Output sample:

Print to stdout, the lowercase version of the sentence, each on a new line.
e.g.

hello codeeval
this is some text
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <algorithm>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    
    while (getline(ifs, line)) {
    	transform(line.begin(), line.end(), line.begin(), ::tolower);
    	cout << line << endl;
    }

    return 0;
}