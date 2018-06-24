/*
Capitalize Words  Share on LinkedIn

Challenge Description:

Write a program which capitalizes words in a sentence.

Input sample:

Your program should accept as its first argument a path to a filename. Input example is the following

Hello world
javaScript language
a letter
Output sample:

Print capitalized words in the following way.

Hello World
JavaScript Language
A Letter
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <string>

using namespace std;

int main( int argc, char** argv ) {
  	ifstream ifs(argv[1]);
  	string line;
  	vector<string> table;

  	while (getline(ifs, line)) {
  		istringstream iss(line);
  		table.clear();

  		string temp;
  		while (iss >> temp) {
  			if ( temp[0] >= 'a' && temp[0] <= 'z' )
  				temp[0] += 'A'-'a';
  			cout << temp << " ";
  		}
  		cout << endl;
  	}

  	return 0;
}