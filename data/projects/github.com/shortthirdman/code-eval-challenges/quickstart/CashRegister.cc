/*
Cash Register  Share on LinkedIn

Description:

The goal of this challenge is to design a cash register program. You will be given two float numbers. The first is the purchase price (PP) of the item. The second is the cash (CH) given by the customer. Your register currently has the following bills/coins within it:

'PENNY': .01,
'NICKEL': .05,
'DIME': .10,
'QUARTER': .25,
'HALF DOLLAR': .50,
'ONE': 1.00,
'TWO': 2.00,
'FIVE': 5.00,
'TEN': 10.00,
'TWENTY': 20.00,
'FIFTY': 50.00,
'ONE HUNDRED': 100.00
The aim of the program is to calculate the change that has to be returned to the customer.
Input sample:

Your program should accept as its first argument a path to a filename.The input file contains several lines. Each line is one test case. Each line contains two numbers which are separated by a semicolon. The first is the Purchase price (PP) and the second is the cash(CH) given by the customer. eg.

15.94;16.00
17;16
35;35
45;50
Output sample:

For each set of input produce a single line of output which is the change to be returned to the customer. In case the CH < PP, print out ERROR. If CH == PP, print out ZERO. For all other cases print the amount that needs to be returned, in terms of the currency values provided. The output should be alphabetically sorted. eg.

NICKEL,PENNY
ERROR
ZERO
FIVE
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    vector<string> table;

    while (getline(ifs, line)) {
    	replace(line.begin(), line.end(), ';', ' ');
    	double price, paid;
        istringstream iss(line);
        iss >> price >> paid;

        table.clear();
        double remain = paid - price;
        if (remain < 0) {
        	cout << "ERROR" << endl;
        	continue;
        } else if (remain == 0) {
            cout << "ZERO" << endl;
        	continue;
        }

        while (1) {
           	if (remain >= 100) {
        		remain -= 100;
                table.push_back("ONE HUNDRED");
        	} else if (remain >= 50) {
        		remain -= 50;
                table.push_back("FIFTY");
        	} else if (remain >= 20) {
        		remain -= 20;
                table.push_back("TWENTY");
        	} else if (remain >= 10) {
        		remain -= 10;
                table.push_back("TEN");
        	} else if (remain >= 5) {
        		remain -= 5;
                table.push_back("FIVE");
        	} else if (remain >= 2) {
        		remain -= 2;
                table.push_back("TWO");
        	} else if (remain >= 1) {
        		remain -= 1;
                table.push_back("ONE");        	
            } else if (remain >= 0.5) {
        		remain -= 0.5;
                table.push_back("HALF DOLLAR");
        	} else if (remain >= 0.25) {
        		remain -= 0.25;
                table.push_back("QUARTER");
        	} else if (remain >= 0.1) {
        		remain -= 0.1;
                table.push_back("DIME");
        	} else if (remain >= 0.049) {
        		remain -= 0.05;
                table.push_back("NICKEL");
        	} else if (remain >= 0.009) {
        		remain -= 0.01;
                table.push_back("PENNY");
        	} else {
                break;
            }
        }

        sort(table.begin(), table.end());
        for (int i = 0; i < table.size(); ++i) {
            if (i > 0) cout << ",";
            cout << table[i];
        }
        cout << endl;
    }

    return 0;
}