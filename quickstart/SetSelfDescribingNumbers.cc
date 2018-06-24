/*
Self Describing Numbers  Share on LinkedIn

Description:

A number is a self-describing number when (assuming digit positions are labeled 0 to N-1), the digit in each position is equal to the number of times that that digit appears in the number.

Input sample:

The first argument is the pathname to a file which contains test data, one test case per line. Each line contains a positive integer. Each line is in the format: N i.e. a positive integer eg.

2020
22
1210
Output sample:

If the number is a self-describing number, print out a 1. If not, print out a 0 eg.

1
0
1
For the curious, here's how 2020 is a self-describing number: Position '0' has value 2 and there is two 0 in the number. Position '1' has value 0 because there are not 1's in the number. Position '2' has value 2 and there is two 2. And the position '3' has value 0 and there are zero 3's.
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    int N;
    string line;
    vector<int> table;
    
    while (getline(ifs, line)) {
        istringstream iss(line);
        iss >> N;

        table.clear();
        table.resize(10);
        
        int digit = 1;
        while (N > digit) {
        	table[(N/digit)%10]++;
        	digit *= 10;
        }
        digit /= 10;

        int pos = 0;
        while (digit > 0) {
        	if (table[pos++] != (N/digit)%10) {
        		cout << "0" << endl;
    			break;
        	} else {
        		digit /= 10;	
        	}
        }
        if (digit == 0) cout << "1" << endl;
    }

    return 0;
}