/*
Lowest Common Ancestor  Share on LinkedIn

Description:

Write a program to determine the lowest common ancestor of two nodes in a binary search tree. You may hardcode the following binary search tree in your program:

    30
    |
  ____
  |   |
  8   52
  |
____
|   |
3  20
    |
   ____
  |   |
  10 29
Input sample:

The first argument will be a text file containing 2 values that represent two nodes within the tree, one per line. e.g. 

8 52
3 29
Output sample:

Print to stdout, the least common ancestor, one per line.
e.g.

30
8
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>

using namespace std;

struct TreeNode {
	int val;
	TreeNode *left;
	TreeNode *right;
	TreeNode(int x) : val(x), left(NULL), right(NULL) {}
};

void generate(TreeNode *root) {
	TreeNode *left1 = new TreeNode(8);
	TreeNode *right1 = new TreeNode(52);
	TreeNode *left2 = new TreeNode(3);
	TreeNode *right2 = new TreeNode(20);
	TreeNode *left3 = new TreeNode(10);
	TreeNode *right3 = new TreeNode(29);

	root->left = left1;
	root->right = right1;
	left1->left = left2;
	left1->right = right2;
	right2->left = left3;
	right2->right = right3;
}

bool exist(TreeNode *root, int val) {
	if (root == NULL) return false;
	if (root->val == val) return true;
	return exist(root->left, val) || exist(root->right, val);
}

int find(TreeNode *root, int A, int B) {
	if (exist(root->left, A)) {
		if (exist(root->left, B)) return find(root->left, A, B);
		else return root->val;
	} else {
		if (exist(root->left, B)) return root->val;
		else return find(root->right, A, B);
	}
}

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    TreeNode *root = new TreeNode(30);
    generate(root);

    while (getline(ifs, line)) {
        istringstream iss(line);
        int A, B;
        iss >> A >> B;
        cout << find(root, A, B) << endl;
    }

    return 0;
}