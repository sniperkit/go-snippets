/*
Sudoku  Share on LinkedIn

Description:

Sudoku is a number-based logic puzzle. It typically comprises of a 9*9 grid with digits so that each column, each row and each of the nine 3*3 sub-grids that compose the grid contains all the digits from 1 to 9. For this challenge, you will be given an N*N grid populated with numbers from 1 through N and you have to determine if it is a valid sudoku solution. You may assume that N will be either 4 or 9. The grid can be divided into square regions of equal size, where the size of a region is equal to the square root of a side of the entire grid. Thus for a 9*9 grid there would be 9 regions of size 3*3 each.

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file contains the value of N, a semicolon and the sqaure matrix of integers in row major form, comma delimited. e.g.

4;1,4,2,3,2,3,1,4,4,2,3,1,3,1,4,2
4;2,1,3,2,3,2,1,4,1,4,2,3,2,3,4,1
Output sample:

Print out True/False if the grid is a valid sudoku layout. e.g.

True
False
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <set>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    vector<vector<int> > table;
    set<int> s;
    string line;
    
    while (getline(ifs, line)) {
    	int range = line[0]-'0';
    	int subrange;
    	if (range == 4) subrange = 2;
    	else subrange = 3;

    	line = line.substr(2);
    	replace(line.begin(), line.end(), ',', ' ');
    	istringstream iss(line);

    	int cur;
    	table.clear();
    	vector<int> x;
    	while (iss >> cur) {
    		x.push_back(cur);
    		if (x.size() == range) {
    			table.push_back(x);
    			x.clear();
    		}
    	}

        // valid row:
    	bool valid = true;
        for (int i = 0; i < table.size(); ++i) {
        	s.clear();
        	for (int j = 0; j < table[0].size(); ++j) {
        		if (table[i][j] < 0 || table[i][j] > range) {
        			valid = false;
        			i = range;
        			j = range;
        		} else {
        			s.insert(table[i][j]);
        		}
        	}
        	if (s.size() != range) { 
        		valid = false;
        		break;
        	}
        }
        if (!valid) {
        	cout << "False" << endl;
        	continue;
        }

        // valid col:
		for (int i = 0; i < table.size(); ++i) {
        	s.clear();
        	for (int j = 0; j < table[0].size(); ++j) {
        		if (table[j][i] < 0 || table[j][i] > range) {
        			valid = false;
        			i = range;
        			j = range;
        		} else {
        			s.insert(table[j][i]);	
        		}
        	}
        	if (s.size() != range) { 
        		valid = false;
        		break;
        	}
        }        
        if (!valid) {
        	cout << "False" << endl;
        	continue;
        }

        // valid square:
        for (int i = 0; i < subrange; ++i) {
        	s.clear();
        	for (int j = 0; j < subrange; ++j) {
        		for (int x = i*subrange; x < i*subrange+subrange; ++x) {
        			for (int y = j*subrange; y < j*subrange+subrange; ++y) {
        				if (table[x][y] < 0 || table[x][y] > range) {
        					valid = false;
        					x = range;
        					y = range;
        				}
        				s.insert(table[x][y]);
        			}
        		}
        	}
        	if (s.size() != range) {
        		valid = false;
        		break;
        	}
        }
        if (!valid) {
        	cout << "False" << endl;
        } else {
        	cout << "True" << endl;
        }
    }

    return 0;
}