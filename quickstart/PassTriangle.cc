/*
Pass Triangle  Share on LinkedIn

Description:

By starting at the top of the triangle and moving to adjacent numbers on the row below, the maximum total from top to bottom is 27.

   5
  9 6
 4 6 8
0 7 1 5
5 + 9 + 6 + 7 = 27
Input sample:

Your program should accept as its first argument a path to a filename. Input example is the following

5
9 6
4 6 8
0 7 1 5
You make also check full input file which will be used for your code evaluation.
Output sample:

The correct output is the maximum sum for the triangle. So for the given example the correct answer would be

27
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    vector<int> pre;
    vector<int> cur;
    string line;

    while (getline(ifs, line)) {
        istringstream iss(line);

        cur.clear();        
        int temp;
        while (iss >> temp) cur.push_back(temp);

        for (int i = 0; i < cur.size(); ++i) {
            if (i < cur.size()-1) {
                if (i == 0) cur[i] += pre[i];
                else cur[i] += max(pre[i], pre[i-1]);
            } else if (i > 0) {
                cur[i] += pre[pre.size()-1];
            }
        }

        swap(pre, cur);
    }

    int res = 0;
    for (int i = 0; i < pre.size(); ++i) {
        res = max(res, pre[i]);
    }
    cout << res << endl;

    return 0;
}