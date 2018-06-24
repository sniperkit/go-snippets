/*
Odd Numbers  Share on LinkedIn

Description:

Print the odd numbers from 1 to 99.

Input sample:

None

Output sample:

Print the odd numbers from 1 to 99, one number per line.
*/

#include <iostream>
#include <fstream>
#include <sstream>

using namespace std;

int main(int argc, char *argv[]) {
	for (int i = 1; i < 100; i += 2) {
		cout << i << endl;
	}

    return 0;
}