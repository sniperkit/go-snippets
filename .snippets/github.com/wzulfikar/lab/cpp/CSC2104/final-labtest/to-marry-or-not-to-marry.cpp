#include <iostream>
#include <vector>
#include <string>
#include <sstream>
#include <iterator>
#include <cstring>
#include <vector>
#include <cstring>
#include <algorithm>

using namespace std;

vector<string> split( string str, char sep = ' ' )
{
	vector<string> splitted;

	istringstream stm(str);
	string token;
	while(std::getline(stm, token, sep)) 
		splitted.push_back(token);
	return splitted;
}

string run(string john, string jane) {
	vector<string> john_things = split(john);
	vector<string> jane_things = split(jane);

	int max = john_things.size();
	if (jane_things.size() > max) {
		max = jane_things.size();
	}

	int matches = 0;

	for (int i = 0; i < john_things.size(); ++i)
	{
		for (int j = 0; j < jane_things.size(); ++j)
		{
			if (john_things[i] == jane_things[j]) {
				matches++;
			}
		}
	}

	// cout << "matches: " << matches << " ds " << max/2;

	if (matches > max/2) {
		return "Forever Dead";
	}
	return "Alive";
}

void test() {
	cout << (run("gaming puzzle sleeping", "chocolate sleeping") == "Alive" ? "✔ " : "✘ ");
	cout << (run("food programming anime travelling", "food shopping anime travelling") == "Forever Dead" ? "✔ " : "✘ ");
	cout << (run("pizza burger pasta spaghetti biriyani", "vegetables") == "Alive" ? "✔ " : "✘ ");
}

int main(){
	// test(); return 0;
	
	int n,j;
	string john, jane;

	cin >> n;
	cin >> john;
	cin >> j;
	cin >> jane;

	cout << run(john, jane);
	return 0;
}