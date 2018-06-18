#include <iostream>
#include <vector>
#include <string>
#include <sstream>
#include <iterator>

using namespace std;

string run(string str) {
	// parse str to args
	istringstream buf(str);
	istream_iterator<std::string> beg(buf), end;
	vector<std::string> tokens(beg, end);
	
	std::string::size_type sz;

	// h, w, l, v
	int args[4];
	
	for (int i = 0; i < tokens.size(); ++i)
	{
		args[i] = std::stoi (tokens[i], &sz);
	}

	/**
	 * KEYWORD: "at least 1000 times bigger"
	 */

	// calc volume
	long long int leon = args[0] * args[1] * args[2];
	long long int pyramid = args[3];

	// cout << "leon: " << leon << " pyramid: " << pyramid;

	if (leon * 1000 <= pyramid) {
		return "YES";
	}

	return "NO";
}

void test() {
	cout << (run("1 1 1 1") == "NO" ? "✔ " : "✘ ");
	cout << (run("1 1 1 1000") == "YES" ? "✔ " : "✘ ");
	cout << (run("1 1 1 10000") == "YES" ? "✔ " : "✘ ");
	cout << (run("10 10 10 1000") == "NO" ? "✔ " : "✘ ");
	cout << (run("100 100 1000 1000") == "NO" ? "✔ " : "✘ ");
}

int main(){
	// test(); return 0;
	
	string a;
	cin >> a;
	cout << run(a);
	return 0;
}