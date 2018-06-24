/*
Closest Pair  Share on LinkedIn

Description:

Credits: Programming Challenges by Steven S. Skiena and Miguel A. Revilla

You will be given the x/y co-ordinates of several locations. You will be laying out 1 cable between two of these locations. In order to minimise the cost, your task is to find the shortest distance between a pair of locations, so that pair can be chosen for the cable installation.

Input sample:

Your program should accept as its first argument a path to a filename.The input file contains several sets of input. Each set of input starts with an integer N (0<=N<=10000), which denotes the number of points in this set. The next N line contains the coordinates of N two-dimensional points. The first of the two numbers denotes the X-coordinate and the latter denotes the Y-coordinate. The input is terminated by a set whose N=0. This set should not be processed. The value of the coordinates will be less than 40000 and non-negative. eg.

5
0 2
6 67
43 71
39 107
189 140
0
Output sample:

For each set of input produce a single line of output containing a floating point number (with four digits after the decimal point) which denotes the distance between the closest two points. If there is no such two points in the input whose distance is less than 10000, print the line INFINITY. eg.

36.2215
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <cmath>
#include <iomanip>

using namespace std;

double dist(pair<long, long> a, pair<long, long> b) {
	double res;
	long long temp = 0;
	temp += pow((double)(a.first-b.first), 2);
	temp += pow((double)(a.second-b.second), 2);

	res = sqrt(temp);

	return res;
}


int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
	long N;
    string line;
    vector<pair<long, long> > table;

    while (ifs >> N) {
    	if (N == 0) break;
    	long x, y;
    	table.clear();
    	for (long i = 0; i < N; ++i) {
    		ifs >> x >> y;
    		table.push_back(make_pair<long, long> (x, y));
    	}

    	double res = dist(table[0], table[1]);
    	for (long i = 0; i < N; ++i) {
    		for (long j = i+1; j < N; ++j) {
    			res = min(res, dist(table[i], table[j]));
    		}
    	}

    	if (res > 10000) cout << "INFINITY" << endl;
    	else cout << setiosflags(ios::fixed) << setw(4) << setprecision(4) << res << endl;
    }

    return 0;
}