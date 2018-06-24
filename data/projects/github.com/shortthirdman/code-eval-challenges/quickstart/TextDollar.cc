/*
Text Dollar  Share on LinkedIn

Description:

Credits: This challenge has been authored by Terence Rudkin

You are given a positive integer number. This represents the sales made that day in your department store. The payables department however, needs this printed out in english. NOTE: The correct spelling of 40 is Forty. (NOT Fourty)

Input sample:

Your program should accept as its first argument a path to a filename.The input file contains several lines. Each line is one test case. Each line contains a positive integer. eg.

3
10
21
466
1234
Output sample:

For each set of input produce a single line of output which is the english textual representation of that integer. The output should be unspaced and in Camelcase. Always assume plural quantities. You can also assume that the numbers are < 1000000000 (1 billion). In case of ambiguities eg. 2200 could be TwoThousandTwoHundredDollars or TwentyTwoHundredDollars, always choose the representation with the larger base i.e. TwoThousandTwoHundredDollars. For the examples shown above, the answer would be:

ThreeDollars
TenDollars
TwentyOneDollars
FourHundredSixtySixDollars
OneThousandTwoHundredThirtyFourDollars
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>

using namespace std;

string processThree(int a) {
    
    static string num1[] = {"", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen"};
    static string num10[] = {"", "Ten", "Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy", "Eighty", "Ninety"};

    if (a == 0) return "";

    string res = "";
    if (a >= 100) {
        res = num1[a/100] + "Hundred";
        a %= 100;
    }

    if (a <= 19) {
        res += num1[a];
    } else {
        res += num10[a/10];
        a %= 10;
        if (a > 0)
            res += num1[a];
    }

    return res;
}

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    int num;
    string units[] = {"", "Thousand", "Million", "Billion"};

    while (ifs >> num) {
        if (num == 1) {
            cout << "OneDollar" << endl;
            continue;
        }

        string res = "Dollars";
        int unitIdx = 0;
        while (num > 0) {
            int rem = num % 1000;
            if (rem != 0) {
                res = processThree(rem) + units[unitIdx] + res;
            }
            num /= 1000;
            unitIdx++;
        }

        cout << res << endl;
    }

    return 0;
}