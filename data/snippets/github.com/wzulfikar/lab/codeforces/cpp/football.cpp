#include<iostream>

using namespace std;

string func(string str) {
	return (str.find("0000000") != -1 || str.find("1111111") != -1) ? "YES" : "NO";
}

void test() {
	cout << (func("001001") == "NO" ? "✔ " : "✘ ");
	cout << (func("1000000001") == "YES" ? "✔ " : "✘ ");
}

int main(){
	test();

	// string s;
	// cin >> s;
	// cout << func(s);
}