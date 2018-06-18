#include<iostream>
#include<string>
#include<sstream>
#include<vector>

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

string run(string dict, string rand) {
	// check args
	// cout << "dict: " << dict << " | rand: " << rand << endl;

	vector<string> in_dict = split(dict);
	vector<string> in_rand = split(rand);

	stringstream aladeen;
	
	for (int i = 0; i < in_rand.size(); ++i)
	{
		for (int j = 0; j < in_dict.size(); ++j)
		{
			if (in_dict[j] == in_rand[i]) {
				in_rand[i] = "Aladeen";
			}
		}
	}

	for (int i = 0; i < in_rand.size(); ++i)
	{
		if (i == in_rand.size() - 1) {
			aladeen << in_rand[i];
		} else {
			aladeen << in_rand[i] << "\n";
		}
	}

	// cout << aladeen.str();

	return aladeen.str();
}

void test() {
	cout << (run("Positive Negetive", "Positive Wadiya Negetive") == "Aladeen\nWadiya\nAladeen" ? "✔ " : "✘ ");
	cout << (run("Easy Hard Impossible", "Aladeen Impossible Easy Hard") == "Aladeen\nAladeen\nAladeen\nAladeen" ? "✔ " : "✘ ");
}

int main(){
	// test(); return 0;

	int dict_n;
	bool first_dict = true;
	cin >> dict_n;
	string dict;
	while (dict_n > 0) {
		string in;
		cin >> in;
		if (first_dict || dict_n == 0) {
			dict += in;
		} else {
			dict += " " + in;
		}
		first_dict = false;
		--dict_n;
	}

	int rand_n;
	bool first_rand = true;
	cin >> rand_n;
	string rand;
	while (rand_n > 0) {
		string in;
		cin >> in;
		if (first_rand || rand_n == 0) {
			rand += in;
		} else {
			rand += " " + in;
		}
		first_rand = false;
		--rand_n;
	}

	cout << run(dict, rand);
	return 0;
}