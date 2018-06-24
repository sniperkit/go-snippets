/*
Overlapping Rectangles  Share on LinkedIn

Description:

Given two axis aligned rectangles A and B, determine if the two overlap.

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file contains 8 comma separated co-ordinates. The co-ordinates are upper left x of A, upper left y of A, lower right x of A, lower right y of A, upper left x of B, upper left y of B, lower right x of B, lower right y of B. e.g.

-3,3,-1,1,1,-1,3,-3
-3,3,-1,1,-2,4,2,2
Output sample:

Print out True or False if A and B intersect. e.g.

False
True
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    int x1, y1, x2, y2, x3, y3, x4, y4;
    string line;

    while (getline(ifs, line)) {
    	replace(line.begin(), line.end(), ',', ' ');
	    istringstream iss(line);
	    iss >> x1 >> y1 >> x2 >> y2 >> x3 >> y3 >> x4 >> y4;

	    if (x3 > x2 || x4 < x1 || y3 < y2 || y4 > y1) {
	    	cout << "False" << endl;
	    } else {
	    	cout << "True" << endl;
	    }
    }

    return 0;
}