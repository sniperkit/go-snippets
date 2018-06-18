#include <iostream>

# include "Queue.h"
# include "Stack.h"

using namespace std;

// Thu, 16 Mar 2017 at 18:52:45 MYT
int main () {
	cout << "Starting Labtest 3..\n\n";

	Stack<int> stack;
	stack.push(1);
	stack.push(2);
	stack.push(3);

	while (!stack.isEmpty()) {
		cout << "Popped " << stack.pop() << " from stack\n";
	}

	cout << "\n";

	Queue<int> queue;
	queue.push(1);
	queue.push(2);
	queue.push(3);

	while (!queue.isEmpty()) {
		cout << "Popped " << queue.pop() << " from queue\n";
	}

	return 0;
}