/*
Set Intersection  Share on LinkedIn

Description:

You are given two sorted list of numbers(ascending order). The lists themselves are comma delimited and the two lists are semicolon delimited. Print out the intersection of these two sets.

Input sample:

File containing two lists of ascending order sorted integers, comma delimited, one per line. e.g. 

1,2,3,4;4,5,6
7,8,9;8,9,10,11,12
Output sample:

Print out the ascending order sorted intersection of the two lists, one per line
e.g.

4
8,9
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <set>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    set<int> table;
    
    while (getline(ifs, line)) {
    	int p = line.find(';');
    	string line1 = line.substr(0, p);
    	replace(line1.begin(), line1.end(), ',', ' ');
		string line2 = line.substr(p+1);
    	replace(line2.begin(), line2.end(), ',', ' ');

    	table.clear();
        istringstream iss1(line1);
        int temp;
        while (iss1 >> temp) table.insert(temp);

        istringstream iss2(line2);
        int flag = 0;
        while (iss2 >> temp) {
        	if (table.find(temp) != table.end()) {
        		if (flag == 0) { cout << temp; flag++; }
        		else cout << "," << temp;
        	}
        }
        cout << endl;
    }

    return 0;
}