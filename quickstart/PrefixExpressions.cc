/*
Prefix expressions  Share on LinkedIn

Description:

You are given a prefix expression. Write a program to evaluate it.
Input sample:

The first argument will be an input file with one prefix expression per line. e.g.

* + 2 3 4
Your program has to read this and insert it into any data structure you like. Traverse that data structure and evaluate the prefix expression. Each token is delimited by a whitespace. You may assume that the only valid operators appearing in test data are '+','*' and '/' 

Output sample:

Print to stdout, the output of the prefix expression, one per line. e.g.

20
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <stack>

using namespace std;

bool isop(char op) {
	return op == '+' || op == '*' || op == '/';
}

int cal(int a, int b, char op) {
	switch (op) {
	case '+':
		return a+b;
		break;
	case '*':
		return a*b;
		break;
	case '/':
		return a/b;
		break;
	}

	return 0;
}

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    stack<char> op;
    stack<int> data;

    while (getline(ifs, line)) {
        istringstream iss(line);
        char cur;
        while (!op.empty()) op.pop();
        while (!data.empty()) data.pop();

        while (iss >> cur) {
        	if (isop(cur)) {
        		op.push(cur);
        		continue;
        	}
            int val = cur-'0';

            while (!op.empty() && !data.empty()) {
                int pre = data.top(); data.pop();
                char cmd = op.top(); op.pop();
                val = cal(val, pre, cmd);
            }

            data.push(val);
        }

        cout << data.top() << endl;
    }

    return 0;
}