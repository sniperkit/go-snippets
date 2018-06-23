#include<iostream>
#include<string>
#include<sstream>

using namespace std;

string run(int x) {
	// (1) y = 2 * x + 1
	// (2) z - 2 * y = x
	// (3) w = x + y - z
	int y = 2 * x + 1;
	int z = x + (2 * y);
	int w = x + y - z;

	stringstream ss;
	ss << y << " " << z << " " << w;
	return ss.str();
}

void test() {
	cout << (run(0) == "1 2 -1" ? "✔ " : "✘ ");
	cout << (run(-999) == "-1997 -4993 1997" ? "✔ " : "✘ ");
	cout << (run(999) == "1999 4997 -1999" ? "✔ " : "✘ ");
}

int main(){
	test(); return 0;
	
	int a;
	cin >> a;
	cout << run(a);
	return 0;
}