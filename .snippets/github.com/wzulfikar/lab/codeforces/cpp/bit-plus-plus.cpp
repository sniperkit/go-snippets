#include<iostream>

using namespace std;

int main(){
	char s[4];
	int n,x=0;
	cin>>n;
	for(;n--;){
		cin>>s;
		if(int(s[1])==43)
			x++;
		else
			x--;
	}
	cout<<x;
}