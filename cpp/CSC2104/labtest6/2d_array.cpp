#include <iostream>

using namespace std;

int main(){
	int n, m;

	cout << "Insert value for 'n': ";
	cin >> n;

	cout << "Insert value for 'm': ";
	cin >> m;

	// declare 2d array of n x n
	int arr[n][n];

	cout << "Size of array is " << n << " x " << n << "\n";

	// initiate arrays to 0
	for (int i = 0; i < n; i++)
	{
		for (int j = 0; j < n; j++)
		{
			arr[i][j] = 0;
		}
	}

	for (int i = 0; i < m; i++)
	{
		int u,v; 

		cout << "Insert value for 'u': ";
		cin >> u;

		cout << "Insert value for 'v': ";
		cin >> v;

		// change the cell
		cout << "Changing value of cell " << u << ":" << v << "\n";
		
		// dunno why but it doesn't crash when 
		// i input u or v with value greater or equal than n.
		// it supposed to crash right? -_
		arr[u][v] = 1;
	}

	cout << "Displaying arrays..\n";
	for (int i = 0; i < n; i++)
	{
		for (int j = 0; j < n; j++)
		{
			cout << arr[i][j] << " ";
		}
		cout << endl;
	}

	return 0;
}