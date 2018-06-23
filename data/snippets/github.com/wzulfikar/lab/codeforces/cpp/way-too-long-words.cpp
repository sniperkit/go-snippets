#include<iostream>
#include<string>
#include<sstream>

using namespace std;

string func(string str) {
	const int LEN = str.size();
	if (LEN <= 10) {
		return str;
	}

	stringstream ss;
	ss << str[0] << LEN-2 << str[LEN-1];

	return ss.str();
}

void test() {
	cout << (func("4") == "" ? "✔ " : "✘ ");
	cout << (func("word") == "word" ? "✔ " : "✘ ");
	cout << (func("localization") == "l10n" ? "✔ " : "✘ ");
	cout << (func("internationalization") == "i18n" ? "✔ " : "✘ ");
}

int main(){
	// test();

	int n;
	string word;

	cin >> n;
	for(int i = 0;i < n;i++) {
		cin >> word;
		cout << func(word) << endl;
	}
}