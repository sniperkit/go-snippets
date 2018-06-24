/*
Multiplication Tables  Share on LinkedIn

Description:

Print out the grade school multiplication table upto 12*12.

Input sample:

None

Output sample:

Print out the table in a matrix like fashion, each number formatted to a width of 4 (The numbers are right-aligned and strip out leadeing/trailing spaces on each line). The first 3 line will look like: 
e.g.

1   2   3   4   5   6   7   8   9  10  11  12
2   4   6   8  10  12  14  16  18  20  22  24
3   6   9  12  15  18  21  24  27  30  33  36
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <iomanip>

using namespace std;

int main(int argc, char *argv[]) {
	for (int i = 1; i <= 12; ++i) {
		for (int j = 1; j <= 12; ++j) {
			cout << setw(4) << i*j;
		}
		cout << endl;
	}

    return 0;
}