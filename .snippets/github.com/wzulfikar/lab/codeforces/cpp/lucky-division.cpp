#include<iostream>

using namespace std;

int main(){
	int in;
	cin >> in;
	cout << (in%4==0 || in%7==0 || in%47==0 || in%74==0 || in%477==0  ?  "YES":"NO");
}