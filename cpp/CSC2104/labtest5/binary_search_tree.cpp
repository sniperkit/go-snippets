#include <iostream>
#include <list>

using namespace std;

// template
template<class T>
class Node {
public:
	T value;
	Node<T>* left;
	Node<T>* right;

	Node(T value) {
		this->value = value;
		this->left = NULL;
		this->right = NULL;
	}
};

template<class F>
class BST {
private:
	Node<F> *root;

	void prefix (Node<F>* current) {
		if (current != NULL) {
			cout << current->value << ' ';
			prefix(current->left);
			prefix(current->right);
		}
	}

	void infix (Node<F>* current) {
		if (current != NULL) {
			infix(current->left);
			cout << current->value << ' ';
			infix(current->right);
		}
	}

	void reverseInfix (Node<F>* current) {
		if (current != NULL) {
			reverseInfix(current->right);
			cout << current->value << ' ';
			reverseInfix(current->left);
		}
	}

	void postfix (Node<F>* current) {
		if (current != NULL) {
			postfix(current->left);
			postfix(current->right);
			cout << current->value << ' ';
		}
	}
public:
	BST(){
		root = NULL;
	}

	void insert (F value) {
		if (root == NULL) {
			root = new Node<F>(value);
		} else {
			Node<F>* temp = root;
			while (temp != NULL) {
				if (value > temp->value) {
					if (temp->right == NULL) {
						temp->right = new Node<F>(value);
						break;
					}
					temp = temp->right;
				} else {
					if (temp->left == NULL) {
						temp->left = new Node<F>(value);
						break;
					}
					temp = temp->left;
				}
			}
		}
	}

	bool find (F value) {
		Node<F>* temp = root;
		while (temp != NULL) {
			if (value == temp->value) {
				return true;
			}
			temp = value > temp->value ? temp->right : temp->left;
		}
		return false;
	}

	void display (int order = 0) {
		switch (order) {
			case 0:
			cout << "Displaying tree in PREFIX mode:\n";
			prefix(root);
			break;

			case 1:
			cout << "Displaying tree in INFIX mode:\n";
			infix(root);
			break;

			case 2:
			cout << "Displaying tree in POSTFIX mode:\n";
			postfix(root);
			break;

			case 3:
			cout << "Displaying tree in REVERSE-INVIX mode:\n";
			reverseInfix(root);
			break;

			default:
			prefix(root);
		}
	}
};

int main() {
	BST<int> tree;
	tree.insert(5);
	tree.insert(10);
	tree.insert(3);
	tree.insert(7);
	tree.insert(10);
	tree.insert(7);

	tree.display(3);

	string findNode = tree.find(11) ? "Found" : "Not found";
	cout << "\nNode 11 is " << findNode;

	string findNode2 = tree.find(10) ? "Found" : "Not found";
	cout << "\nNode 10 is " << findNode2;
	return 0;
}
