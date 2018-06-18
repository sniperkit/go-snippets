#include<iostream>
#include<string>
using namespace std;

string func(string str) {
	const string NEEDLE = "WUB";
	const string REPLACEMENT = " ";

	int pos = str.find(NEEDLE);
	while (pos != -1) {
		str.replace(pos, 3, REPLACEMENT);
		pos = str.find(NEEDLE);
	}
	return str;
}

void test() {
	cout << (func("WUBWUBABCWUB") == "ABC" ? "✔ " : "✘");
	cout << (func("WUBWEWUBAREWUBWUBTHEWUBCHAMPIONSWUBMYWUBFRIENDWUB") == "WE ARE  THE CHAMPIONS MY FRIEND" ? "✔ " : "✘");
}

int main(){
	// test();

	string a;
	cin>>a;
	cout << func(a);
}