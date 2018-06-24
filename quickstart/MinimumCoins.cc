/*
Minimum Coins  Share on LinkedIn

Description:

You are given 3 coins of value 1, 3 and 5. You are also given a total which you have to arrive at. Find the minimum number of coins to arrive at it.

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file contains a positive integer number which represents the total you have to arrive at e.g.

11
20
Output sample:

Print out the minimum number of coins required to arrive at the total e.g.

3
4

*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    int N;

    while (ifs >> N) {
    	int count = 0;
    	while (N > 0) {
    		if (N >= 5) {
    			count += N/5;
    			N %= 5;
    		} else if (N >= 3) {
    			count += N/3;
    			N %= 3;
    		} else if (N >= 1) {
    			count += N;
    			break;
    		}
    	}
    	cout << count << endl;
    }

    return 0;
}