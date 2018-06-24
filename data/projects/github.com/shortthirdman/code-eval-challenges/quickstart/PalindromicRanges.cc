/*
Palindromic Ranges  Share on LinkedIn

Description:

A positive integer is a palindrome if its decimal representation (without leading zeros) is a palindromic string (a string that reads the same forwards and backwards). For example, the numbers 5, 77, 363, 4884, 11111, 12121 and 349943 are palindromes.

A range of integers is interesting if it contains an even number of palindromes. The range [L, R], with L <= R, is defined as the sequence of integers from L to R (inclusive): (L, L+1, L+2, ..., R-1, R). L and R are the range's first and last numbers.

The range [L1,R1] is a subrange of [L,R] if L <= L1 <= R1 <= R. Your job is to determine how many interesting subranges of [L,R] there are.

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file is one test case. Each test case will contain two positive integers, L and R (in that order), separated by a space. eg.

1 2
1 7
87 88
Output sample:

For each line of input, print out the number of interesting subranges of [L,R] eg.

1
12
1
For the curious: In the third example, the subranges are: [87](0 palindromes), [87,88](1 palindrome),[88](1 palindrome). Hence the number of interesting palindromic ranges is 1
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <set>

using namespace std;

bool ispalin(int i) {
    stringstream ss;
    ss << i;

    string s = ss.str();
    string reversed = s;
    reverse(reversed.begin(), reversed.end());

    return (s == reversed);
}


int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    int low, high;

    while (ifs >> low >> high) {
        set<int> palin;
        int tmp = low;
        while (tmp <= high) {
            if (ispalin(tmp))
                palin.insert(tmp);
            tmp++;
        }

        int res = 0;
        for (int i = low; i <= high; ++i) {
            bool pre = false;
            if (palin.find(i) == palin.end()) {
                pre = true;
                ++res;
            }

            for (int j = i+1; j <= high; ++j) {
                if (palin.find(j) == palin.end()) {
                    if (pre == true) res++;
                } else {
                    if (pre == true) {
                        pre = false;
                    } else {
                        pre = true;
                        res++;
                    }
                }
            }
        }

        cout << res << endl;
    }

    return 0;
}