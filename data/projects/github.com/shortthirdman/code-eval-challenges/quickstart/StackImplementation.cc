/*
Stack Implementation  Share on LinkedIn

Description:

Write a program implementing a stack inteface for integers.The interface should have 'push' and 'pop' functions. You will be asked to 'push' a series of integers and then 'pop' and print out every alternate integer.
Input sample:

The first argument will be a text file containing a series of space delimited integers, one per line. e.g. 

1 2 3 4
10 -2 3 4
Output sample:

Print to stdout, every alternate integer(space delimited), one per line.
e.g.

4 2
4 -2
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    vector<int> table;
    string line;
    
    while (getline(ifs, line)) {
        istringstream iss(line);
        table.clear();
        int temp;
        while (iss >> temp) table.push_back(temp);

        for (int i = table.size()-1; i >= 0; i -= 2) {
        	cout << table[i] << " ";
        }
        cout << endl;
    }

    return 0;
}