/*
Rightmost Char  Share on LinkedIn

Description:

You are given a string 'S' and a character 't'. Print out the position of the rightmost occurrence of 't'(case matters) in 'S' or -1 if there is none. The position to be printed out is zero based.

Input sample:

The first argument is a file, containing a string and a character, comma delimited, one per line. Ignore all empty lines in the input file.e.g. 

Hello World,r
Hello CodeEval,E
Output sample:

Print out the zero based position of the character 't' in string 'S', one per line. Do NOT print out empty lines between your output.
e.g.

8
10
*/

#include <iostream>
#include <fstream>
#include <sstream>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    
    while (getline(ifs, line)) {
    	int p = line.find(',');
    	string left = line.substr(0, p);
    	char t = line[line.size()-1];
    	string::size_type res = left.find_last_of(t);
    	if (res != string::npos)
	    	cout << res << endl;
	    else
	    	cout << "-1" << endl;
    }

    return 0;
}