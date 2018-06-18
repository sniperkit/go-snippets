#include <iostream>

using namespace std;

template<class T>
class Node{
	public:
		T value;
		Node<T>* left;
		Node<T>* right;

		Node(T value){
			this->value = value;
			this->left = NULL;
			this->right = NULL;
		}
};

template<class F>
class BST{ // Binary Search Tree (BST)
	private:
		Node<F>* root;
	public:
		BST(){
			root = NULL;
		}

		void insert(F value){
			if(root==NULL){
				root = new Node<F>(value);
			}else{
				Node<F>* temp = root;
				while(temp!=NULL){
					if(value > temp->value){ 
					// go right, new value is greater than current tree node:w
						if(temp->right == NULL){
							temp->right = new Node<F>(value);
							break;
						}
						temp = temp->right;
					}else{ // go left, new value is less than current tree node
						if(temp->left == NULL){
							temp->left = new Node<F>(value);
							break;
						}
						temp = temp->left;
					}
				}
			}
		}

		bool find(F value){
			Node<F>* temp = root;
			while(temp!=NULL){
				if(value > temp->value){ 
					temp = temp->right;
				}else if(value < temp->value){ 
					temp = temp->left;
				}else{
					return true;
				}
			}
			return false;
		}

		Node<F>* getMinNode(Node<F>* node){
			while(node->left!=NULL){
				node = node->left;
			}
			return node;
		}

		void deleteHelper(F value, Node<F>* current, Node<F>* parent = NULL){
			if(value < current->value){
			
				// 'value' is smaller, going left to search further
				deleteHelper(value, current->left, current);
			
			}else if(value > current->value){
			
				// 'value' is bigger, going right to search further
				deleteHelper(value, current->right, current);
			
			}else if(value == current->value){ 
				// Found match for the 'value'

				if(current->left!=NULL && current->right!=NULL){
					
					// Has two child node
					Node<F>* minNode = getMinNode(current->right);
					current->value = minNode->value;
					deleteHelper(minNode->value, current->right, current);
					
				}else if(current->left!=NULL){
					
					// Has only one child node, left child
					if(parent!=NULL){
						if(parent->left == current) parent->left = current->left; // LL
						else parent->right = current->left; // RL
					}else{
						root = current->left;
					}
					
					delete current;

				}else if(current->right!=NULL){
					
					// Has only one child node, right child
					if(parent!=NULL){
						if(parent->left == current) parent->left = current->right; // LR
						else parent->right = current->right; // RR
					}else{
						root = current->right;
					}
					
					delete current;
				
				}else{
					
					// Has no child node
					if(parent!=NULL){
						if(parent->left == current) parent->left = NULL;
						else parent->right = NULL;
					}else{
						root = NULL;
					}
					
					delete current;
				}
			}
		}

		void deleteNode(F value){	
			if(root!=NULL){
				deleteHelper(value, root);
			}
		}

		void prefix(Node<F>* current){
			if(current!=NULL){
				cout<<current->value<<" ";
				prefix(current->left);
				prefix(current->right);
			}
		}

		void infix(Node<F>* current){
			if(current!=NULL){
				infix(current->left);
				cout<<current->value<<" ";
				infix(current->right);
			}
		}

		void postfix(Node<F>* current){
			if(current!=NULL){
				postfix(current->left);
				postfix(current->right);
				cout<<current->value<<" ";
			}
		}

		void display(int order=0){
			if(order==0){
				prefix(root);
			}else if(order==1){
				infix(root);
			}else if(order==2){
				postfix(root);
			}
			cout<<endl;
		}
};

int main(){
	BST<int> tree;
	int nums[] = {100,50,200,150,250,300};
	
	for(int i=0;i<6;i++){
		tree.insert(nums[i]);
	}

	tree.display(1); // Correct Output: 300 250 200 150 100 50 
	tree.deleteNode(200);
	tree.display(); // Correct Output: 100 150 250 300 50
	cout<<tree.find(200); // Correct Output: 0
	cout<<endl;
	cout<<tree.find(300); // Correct Output: 1

	return 0;
}
