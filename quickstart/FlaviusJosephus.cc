/*
Flavius Josephus  Share on LinkedIn

Description:

Flavius Josephus was a famous Jewish historian of the first century, at the time of the destruction of the Second Temple. According to legend, during the Jewish-Roman war he was trapped in a cave with a group of soldiers surrounded by Romans. Preferring death to capture, the Jews decided to form a circle and, proceeding around it, to kill every j'th person remaining until no one was left. Josephus found the safe spot in the circle and thus stayed alive.Write a program that returns a list of n people, numbered from 0 to n-1, in the order in which they are executed.

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file contains two comma separated positive integers n and m , where n is the number of people and every m'th person will be executed. e.g.

10,3
5,2
Output sample:

Print out the list of n people(space delimited) in the order in which they will be executed. e.g.

2 5 8 1 6 0 7 4 9 3
1 3 0 4 2
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    int N, k;
   	string line;
    vector<int> table;

    while (getline(ifs, line)) {
    	replace(line.begin(), line.end(), ',', ' ');
        istringstream iss(line);
        iss >> N >> k;

        for (int i = 0; i < N; ++i) table.push_back(i);
        int step = 1;
        int cur = 0;
        int remain = N;
        while (1) {
            // advance to next available element
            while (table[cur] == -1)
                if (++cur == table.size()) cur = 0;

            // tick one out
        	if (step % k == 0) {
        		cout << table[cur] << " ";
                table[cur] = -1;
                if (--remain == 0) {
                    cout << endl;
                    break;
                }
        	}
            ++step;
            if (++cur == table.size()) cur = 0;
        }
    }

    return 0;
}