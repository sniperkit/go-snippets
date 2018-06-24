/*
Following Integer  Share on LinkedIn

Description:

Credits: This challenge has appeared in a past google competition

You are writing out a list of numbers.Your list contains all numbers with exactly Di digits in its decimal representation which are equal to i, for each i between 1 and 9, inclusive. You are writing them out in ascending order. For example, you might be writing every number with two '1's and one '5'. Your list would begin 115, 151, 511, 1015, 1051. Given N, the last number you wrote, compute what the next number in the list will be. The number of 1s, 2s, ..., 9s is fixed but the number of 0s is arbitrary.

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file is one test case. Each test case will contain an integer n < 10^6

Output sample:

For each line of input, generate a line of output which is the next integer in the list.
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    vector<int> data;
    string line;
    
    while (getline(ifs, line)) {
        data.clear();
        for (int i = 0; i < line.size(); ++i) 
        	data.push_back(line[i]-'0');

        int cur = data.size()-2;
        while (cur >= 0 && data[cur] >= data[cur+1]) --cur;
        
        if (cur == -1) {
            reverse(data.begin(), data.end());
            if (data[0] == 0) {
                int idx = 1;
                while (data[idx] == 0) idx++;
                swap(data[0], data[idx]);
            }

        	for (int i = 0; i < data.size(); ++i) {
        		cout << data[i];
        		if (i == 0) cout << "0";
        	}
        	cout << endl;
        } else {
        	int j = data.size()-1;
        	while (j > cur && data[j] <= data[cur]) --j;
        	swap(data[cur], data[j]);
        	reverse(data.begin()+cur+1, data.end());
        	for (int i = 0; i < data.size(); ++i)
        		cout << data[i];
        	cout << endl;
        }
    }

    return 0;
}