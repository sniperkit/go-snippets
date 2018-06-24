/*
Endianness  Share on LinkedIn

Description:

Write a program to print out the endianness of the system.

Input sample:

None

Output sample:

Print to stdout, the endianness, wheather it is LittleEndian or BigEndian.
e.g.

BigEndian
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>

using namespace std;

int main(int argc, char *argv[]) {
	short int word = 0x0001;
	char *byte = (char *)&word;
	if (byte[0] == 0) {
		cout << "BigEndian" << endl;
	} else {
		cout << "LittleEndian" << endl;
	}

    return 0;
}