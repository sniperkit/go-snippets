#include <iostream>

using namespace std;

template<class K, class V>
class Node{
	public:
		K key;
		V value;
		int height;
		Node<K,V>* left;
		Node<K,V>* right;

		Node(K key, V value){
			this->key = key;
			this->value = value;
			this->height = 1;
			this->left = NULL;
			this->right = NULL;
		}
};

template<class K, class V>
class AVLTree{
	private:
		Node<K,V>* root;
	public:
		AVLTree(){
			root = NULL;
		}

		int get_height(Node<K,V>* cur){
			int lh = (cur->left==NULL ? 0:cur->left->height);
			int rh = (cur->right==NULL ? 0:cur->right->height);
			return max(lh,rh)+1;
		}

		int get_balance_factor(Node<K,V>* node){
			int b = (node->left==NULL ? 0:node->left->height) - (node->right==NULL ? 0:node->right->height);
			return b;
		}

		Node<K,V>* ll_rotation(Node<K,V>* z){
			Node<K,V>* y = z->left;
			z->left = y->right;
			y->right = z;
			z->height-=2;
			return y;
		}

		Node<K,V>* rr_rotation(Node<K,V>* z){
			Node<K,V>* y = z->right;
			z->right = y->left;
			y->left = z;
			z->height-=2;
			return y;
		}

		Node<K,V>* lr_rotation(Node<K,V>* z){
			Node<K,V>* y = z->left;
			z->left = ll_rotation(y);
			return rr_rotation(z);
		}

		Node<K,V>* rl_rotation(Node<K,V>* z){
			Node<K,V>* y = z->right;
			z->right = rr_rotation(y);
			return ll_rotation(z);
		}

		void insert_helper(K key, V value, Node<K,V>*& current, bool& isAdded){
			if(current==NULL){
				current = new Node<K,V>(key, value);
				isAdded = true;
				return;
			}
			if(key < current->key){
				insert_helper(key, value, current->left, isAdded);
			}else if(key > current->key){
				insert_helper(key, value, current->right, isAdded);
			}else{
				// the 'key' already exists in the tree
				return;
			}

			if(isAdded){
				current->height = get_height(current);

				int b = get_balance_factor(current);

				if(b < -1){
					if(get_balance_factor(current->right) > 0){
						current = rl_rotation(current);
					}else{
						current = rr_rotation(current);
					}
				}else if(b > 1){
					if(get_balance_factor(current->left) > 0){
						current = ll_rotation(current);
					}else{
						current = lr_rotation(current);
					}
				}
			}

		}

		void insert(K key, V value){
			bool isAdded = false;
			insert_helper(key, value, root, isAdded);	
		}

		void preOrder(Node<K,V>* current, int level=0){
			if(current!=NULL){
				cout<<"("<<current->key<<","<<current->value<<") ";
				preOrder(current->left,level+1);
				preOrder(current->right,level+1);
			}
		}
		
		void display(int order=0){
			if(order==0){
				preOrder(root);
			}
			cout<<endl;
		}

};

int main(){
	AVLTree<int, char> tree;

	for(int i=0;i<10;i++){
		char ch = 'A'+i;
		tree.insert(i+1,ch);
		tree.display();
	}

	return 0;
}