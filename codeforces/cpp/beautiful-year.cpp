#include<iostream>
using namespace std;

main()
{
	int a,a1=0,a2=0,a3=0,a4=0;
	cin>>a;
	while(a1==a2||a1==a3||a1==a4||a2==a3||a2==a4||a3==a4)
	{
		a=a+1;
		a1=a%10;
		a2=a/10%10;
		a3=a/100%10;
		a4=a/1000%10;
	}
	cout<<a;
}