/*
Simple Sorting  Share on LinkedIn


Challenge Description:

Write a program which sorts numbers.

Input sample:

Your program should accept as its first argument a path to a filename. Input example is the following

70.920 -38.797 14.354 99.323 90.374 7.581
-37.507 -3.263 40.079 27.999 65.213 -55.552
Output sample:

Print sorted numbers in the following way.

-38.797 7.581 14.354 70.920 90.374 99.323
-55.552 -37.507 -3.263 27.999 40.079 65.213
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <algorithm>
#include <iomanip>
#include <vector>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    vector<double> data;
    
    while (getline(ifs, line)) {
        istringstream iss(line);
        data.clear();
        double tmp;
        while (iss >> tmp) {
        	data.push_back(tmp);
        }
        sort(data.begin(), data.end());

        for (int i = 0; i < data.size(); ++i) {
        	if (i > 0) cout << " ";
        	cout << fixed << setprecision(3) << data[i];
        }
        cout << endl;
    }

    return 0;
}