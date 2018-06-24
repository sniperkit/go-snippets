/*
Grid Walk  Share on LinkedIn

Description:

There is a monkey which can walk around on a planar grid. The monkey can move one space at a time left, right, up or down. That is, from (x, y) the monkey can go to (x+1, y), (x-1, y), (x, y+1), and (x, y-1). Points where the sum of the digits of the absolute value of the x coordinate plus the sum of the digits of the absolute value of the y coordinate are lesser than or equal to 19 are accessible to the monkey. For example, the point (59, 79) is inaccessible because 5 + 9 + 7 + 9 = 30, which is greater than 19. Another example: the point (-5, -7) is accessible because abs(-5) + abs(-7) = 5 + 7 = 12, which is less than 19. How many points can the monkey access if it starts at (0, 0), including (0, 0) itself?

Input sample:

There is no input for this program.
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <set>

using namespace std;

bool isvalid(long x, long y) {
	x = abs(x);
	y = abs(y);

	int res = 0;

	while (x > 0) {
		res += x % 10;
		x /= 10;
	}
	if (res > 19) return false;

	while (y > 0) {
		res += y % 10;
		y /= 10;
	}
	if (res > 19) return false;

	return true;
}

void process(vector<pair<long, long> > &table, set<pair<long, long> > &set, long x, long y) {
	if (isvalid(x, y) && set.find(make_pair<long, long> (x, y)) == set.end()) {
		set.insert(make_pair<long, long> (x, y));
		table.push_back(make_pair<long, long> (x, y));
	}
}

int main(int argc, char *argv[]) {
    vector<pair<long, long> > table;
    set<pair<long, long> > set;
    table.push_back(make_pair<long, long> (0, 0));
    set.insert(make_pair<long, long> (0, 0));

    for (long i = 0; i < table.size(); ++i) {
    	long x = table[i].first;
    	long y = table[i].second;

    	process(table, set, x-1, y-1);
    	process(table, set, x+1, y+1);
    	process(table, set, x+1, y-1);
    	process(table, set, x-1, y+1);
    	process(table, set, x-1, y);
    	process(table, set, x, y-1);
    	process(table, set, x+1, y);
    	process(table, set, x, y+1);
    }

    cout << table.size() << endl;

    return 0;
}