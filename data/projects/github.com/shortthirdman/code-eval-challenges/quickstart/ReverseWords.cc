/*
Reverse words  Share on LinkedIn

Description:

Write a program to reverse the words of an input sentence.
Input sample:

The first argument will be a text file containing multiple sentences, one per line. Possibly empty lines too. e.g. 

Hello World
Hello CodeEval
Output sample:

Print to stdout, each line with its words reversed, one per line. Empty lines in the input should be ignored. Ensure that there are no trailing empty spaces on each line you print. 
e.g.

World Hello
CodeEval Hello
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    vector<string> table;
    
    while (getline(ifs, line)) {
        istringstream iss(line);
        table.clear();
        
        string temp;
        while (iss >> temp) 
        	table.push_back(temp);
        
        for (int i = table.size()-1; i >= 0; --i) 
        	cout << table[i] << " ";
        cout << endl;
    }

    return 0;
}