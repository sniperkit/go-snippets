/*
Penultimate Word  Share on LinkedIn

Challenge Description:

Write a program finds the next-to-last word in a string.

Input sample:

Your program should accept as its first argument a path to a filename. Input example is the following

some line with text
another line
Each line has more than one word.

Output sample:

Print the next-to-last word in the following way.

with
another
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <algorithm>
#include <vector>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    vector<string> data;
    
    while (getline(ifs, line)) {
        istringstream iss(line);
        data.clear();
        string tmp;
        while (iss >> tmp) {
        	data.push_back(tmp);
        }

        if (data.size() >= 2) 
            cout << data[data.size()-2] << endl;
    }

    return 0;
}