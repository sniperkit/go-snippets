/*
Ugly Numbers  Share on LinkedIn

Description:

Credits: This challenge has appeared in a google competition before.

Once upon a time in a strange situation, people called a number ugly if it was divisible by any of the one-digit primes (2, 3, 5 or 7). Thus, 14 is ugly, but 13 is fine. 39 is ugly, but 121 is not. Note that 0 is ugly. Also note that negative numbers can also be ugly; -14 and -39 are examples of such numbers.

One day on your free time, you are gazing at a string of digits, something like:

123456
You are amused by how many possibilities there are if you are allowed to insert plus or minus signs between the digits. For example you can make:
1 + 234 - 5 + 6 = 236
which is ugly. Or
123 + 4 - 56 = 71
which is not ugly.

It is easy to count the number of different ways you can play with the digits: Between each two adjacent digits you may choose put a plus sign, a minus sign, or nothing. Therefore, if you start with D digits there are 3^(D-1) expressions you can make. Note that it is fine to have leading zeros for a number. If the string is '01023', then '01023', '0+1-02+3' and '01-023' are legal expressions.

Your task is simple: Among the 3^(D-1) expressions, count how many of them evaluate to an ugly number.
Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file is one test case. Each test case will be a single line containing a non-empty string of decimal digits. The string in each test case will be non-empty and will contain only characters '0' through '9'. Each string is no more than 13 characters long. eg.

1
9
011
12345
Output sample:

Print out the number of expressions that evaluate to an ugly number for each test case, each one on a new line eg

0
1
6
64
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>

using namespace std;

void generate(vector<string> &table, string line, int i) {
	if (i == line.size()) {
		table.push_back(line);
		return;
	}

	string line1 = line.substr(0,i) + "+" + line.substr(i);
	string line2 = line.substr(0,i) + "-" + line.substr(i);

	generate(table, line1, i+2);
	generate(table, line2, i+2);
	generate(table, line, i+1);
}

// use a default val. For example, 1+2 is treated as 0 + 1+2
bool ugly(string line) {
	long pre = 0;
	long cur = 0;
	char op = '+';

	for (int i = 0; i < line.size(); ++i) {
		if (line[i] != '+' && line[i] != '-') {
			cur *= 10;
			cur += line[i]-'0';
		} else {
			if (op == '+')	cur = pre + cur;
			else			cur = pre - cur;
			pre = cur;
			cur = 0;
			op = line[i];
		}
	}

	if (op == '+')	cur = pre + cur;
	else			cur = pre - cur;

	// cout << line << " = " << cur << endl;

	return (cur%2 == 0) || (cur%3 == 0) || (cur%5 == 0) || (cur%7 == 0);
}


int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    vector<string> table;
    
    while (getline(ifs, line)) {
    	table.clear();
    	generate(table, line, 1);

    	int res = 0;
    	for (int i = 0; i < table.size(); ++i) {
    		if (ugly(table[i])) ++res;
    	}

    	cout << res << endl;
    }

    return 0;
}