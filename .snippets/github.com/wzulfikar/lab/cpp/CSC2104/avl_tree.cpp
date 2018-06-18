#include <iostream>

using namespace std;

template<class K, class V>
class Node {
public:
	K key;
	V value;
	int height;
	Node<K, V>* left;
	Node<K, V>* right;

	Node (K key, V value) {
		this->key = key;
		this->value = value;
		this->height = 1;
		this->left = NULL;
		this->right = NULL;
	}
};

template<class K, class V>
class AVLTree {
private:
	Node<K, V>* root;
public:
	AVLTree () {
		root = NULL;
	}

	int getBalanceFactor (Node<K, V>* node) {
		int b = getHeight(node->left) - getHeight(node->right);
		return b;
	}

	Node<K, V>* ll_rotation (Node<K,V>* z) {
		
	}

	Node<K, V>* rr_rotation (Node<K,V>* z) {

	}

	Node<K, V>* lr_rotation (Node<K,V>* z) {
		
	}

	Node<K, V>* rl_rotation (Node<K,V>* z) {

	}

	void insertHelper (K key, V value, Node<K,V>*& current, bool& isAdded) {
		if (current == NULL) {
			current = new Node<K,V>(key, value);
			isAdded = true;
			return;
		}
		if (key < current->key) {
			insertHelper(key, value, current->left, isAdded);
		} else if (key > current->key) {
			insertHelper(key, value, current->right, isAdded);
		} else {
			// the 'key' is already exists
			return;
		}
	}

	void insert (K key, V value) {
		bool isAdded = false;
		insertHelper(key, value, root, isAdded);
	}

	void preOrder (Node<K, V>* cur) {
		if (cur != NULL) {
			cout << " (" << cur->key << ", " << cur->value << ") ";
			preOrder(cur->left);
			preOrder(cur->right);
		}
	}

	void display (int order = 0) {
		if (order == 0) {
			preOrder(root);
		}
		cout << endl;
	}
};

int main() {
	return 0;
}
