/*
Mth to last element  Share on LinkedIn

Description:

Write a program to determine the Mth to last element of a list.
Input sample:

The first argument will be a text file containing a series of space delimited characters followed by an integer representing a index into the list(1 based), one per line. e.g. 

a b c d 4
e f g h 2
Output sample:

Print to stdout, the Mth element from the end of the list, one per line. If the index is larger than the list size, ignore that input. 
e.g.

a
g
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    vector<char> table;
    string line;
    
    while (getline(ifs, line)) {
        istringstream iss(line);
        char temp;
        table.clear();
        while (iss >> temp) table.push_back(temp);

        int pos = table.back()-'0';
        cout << table[table.size()-1-pos] << endl;
    }

    return 0;
}