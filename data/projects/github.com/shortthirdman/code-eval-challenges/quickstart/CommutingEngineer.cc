/*
Commuting Engineer

Challenge Description:

We spend a lot of time and money commuting. As well as many times we need to stop by at several places in the City. Let's stop commuting and live in San Francisco!

Solving the following challenge means finding the shortest possible route which visits each coordinate once starting from the point 1. You may read more about the Travelling salesman problem

On the map you can see the best route for 5 coordinates. But we've added 5 more. So you need to find the best route for 10 coordinates.

Input sample:

Your program should accept as its first argument a path to a filename. Input example is the following

1 | 1299 El Camino Del Mar, SF (37.7842044, -122.50328830000001)
2 | 101 Cervantes Boulevard, SF (37.8041642, -122.43961819999998)
3 | 1120 Noe Street, SF (37.7511053, -122.43185189999997)
4 | 1201 Funston Avenue, SF (37.765457, -122.47103800000002)
5 | 1000 Lombard Street, SF (37.802458, -122.418207)
Output sample:

It must start from position 1

1
4
3
5
2

*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <cstdio>
#include <iomanip>
#include <cmath>
#include <cstdlib>
#include <algorithm>

using namespace std;

double stringToNum(string s) {
	double res = 0;
	bool flag = false;
	if (s[0] == '-') {
		flag = true;
		s.substr(1);
	}

	res = atof(s.c_str());
	if (flag) res = -res;
	return res;
}

void parse(vector<pair<double, double> > &points, string buf) {
	buf = buf.substr(buf.find('(')+1);	// get rid of first parts
	string left = buf.substr(0, buf.find(','));	// get left number
	buf = buf.substr(buf.find(',')+2);	// get rid of ", "
	string right = buf.substr(0, buf.size()-1);	// get rid of ')'

	points.push_back(make_pair<double, double> (stringToNum(left), stringToNum(right)));
}

double dist(pair<double, double> &a, pair<double, double> &b) {
	return sqrt(pow(a.first-b.first, 2) + pow(a.second-b.second, 2));
}

double route(vector<pair<double, double> > &points, vector<int> &order) {
	double res = 0;
	for (int i = 1; i < order.size(); ++i) {
		res += dist(points[order[i-1]], points[order[i]]);
	}
	return res;
}

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string buf;
    vector<pair<double, double> > points;
    vector<int> order;

    while (getline(ifs, buf)) {
    	if (buf.size() == 0) continue;
    	parse(points, buf);
    	order.push_back(order.size());
    }

    double res = 9999.9999;
    vector<int> resOrder;

    do {
    	double temp = route(points, order);

    	if (temp < res) {
    		res = temp;
    		resOrder = order;
    	}
    } while (next_permutation(order.begin()+1, order.end()));


    for (int i = 0; i < resOrder.size(); ++i) {
    	cout << resOrder[i]+1 << endl;
    }

    return 0;
}