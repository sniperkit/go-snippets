# include "DList.h"

int randomInt (int MAX = 100) {
	return (rand() % MAX) + 1;
}

void makeListOfTenNumbers (DList<int> list) {
	cout << "Making list of 10 numbers..\n";

	int N = 10;
	for (int i = 0; i < N; i++) {
		int number = randomInt();
		list.push_back(number);
	}
}

void deleteList (DList<int> list) {
	cout << "Deleting list..\n";
	int deletedItem;
	while (!list.isEmpty()) {
		list.delete_front();
		cout << "Deleted " << deletedItem << " from list\n";
	}
}

int recursive (int N) {
	return N <= 2 ? N : recursive(N - 3) + recursive(N - 2) + recursive(N - 1);
}

int main () {

	cout << "Starting Labtest 2..\n\n";

	DList<int> list;
	makeListOfTenNumbers(list);
	list.display();

	cout << "Printing in reverse order..\n";
	list.display(true);

	deleteList(list);
	list.display();
	return 0;
}
