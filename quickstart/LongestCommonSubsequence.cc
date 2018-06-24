/*
Longest Common Subsequence  Share on LinkedIn

Description:

You are given two sequences. Write a program to determine the longest common subsequence between the two strings(each string can have a maximum length of 50 characters). (NOTE: This subsequence need not be contiguous. The input file may contain empty lines, these need to be ignored.)
Input sample:

The first argument will be a file that contains two strings per line, semicolon delimited. You can assume that there is only one unique subsequence per test case. e.g.

XMJYAUZ;MZJAWXU
Output sample:

The longest common subsequence. Ensure that there are no trailing empty spaces on each line you print. e.g.

MJAU
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <stack>
#include <iomanip>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;

    while (getline(ifs, line)) {
        if (line.size() == 0) continue;
        int pos = line.find(';');
        string a = line.substr(0, pos);
        string b = line.substr(pos+1);

        int m[a.size()+1][b.size()+1];

        for (int i = 0; i <= a.size(); ++i)
            m[i][0] = 0;
        for (int i = 0; i <= b.size(); ++i)
            m[0][i] = 0;

        for (int i = 1; i <= a.size(); ++i) {
            for (int j = 1; j <= b.size(); ++j) {
                if (a[i-1] == b[j-1])
                    m[i][j] = m[i-1][j-1]+1;
                else
                    m[i][j] = max(m[i][j-1], m[i-1][j]);
            }
        }

        int i = a.size();
        int j = b.size();
        stack<char> s;

        while (i > 0 && j > 0) {
            if (a[i-1] == b[j-1]) {
                s.push(a[i-1]);
                --i;  --j;
            } else {
                if (m[i][j] == m[i-1][j]) --i;
                else                      --j;
            }
        }

        while (!s.empty()) {
            cout << s.top();  s.pop();
        }
        cout << endl;
    }

    return 0;
}