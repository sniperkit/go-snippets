/*
Type Ahead  Share on LinkedIn

Description:

Your task is to build a type-ahead feature for an upcoming product.

When the user enters a word or set of words, we want to be able to "predict" what they're going to type next with some level of accuracy. We've chosen to implement this using the N-Gram algorithm defined here http://en.wikipedia.org/wiki/N-gram.

Your program should return a tuple of predictions sorted high to low based on the prediction score (upto a maximum of three decimal places, or pad with zeroes upto three decimal places i.e. 0.2 should be shown as 0.200), (if predictions share the same score, they are sorted alphabetically.) Words should be split by whitespace with all non-alphanumeric characters stripped off the beginning and end.

Prediction scores are calculated like this:
Occurrences of a word after an N-gram / total number of words after an N-gram
e.g. for an N-gram of length 2:

ONE TWO ONE TWO THREE TWO THREE
"TWO" has the following predictions:
THREE:0.666,ONE:0.333
"THREE" occurred 2 times after a "TWO" and "ONE" occurred 1 time after a "TWO", for a total of 3 occurrences after a "TWO".

Your program will run against the following text i.e. Harcode it into your program:
Mary had a little lamb its fleece was white as snow;
And everywhere that Mary went, the lamb was sure to go. 
It followed her to school one day, which was against the rule;
It made the children laugh and play, to see a lamb at school.
And so the teacher turned it out, but still it lingered near,
And waited patiently about till Mary did appear.
"Why does the lamb love Mary so?" the eager children cry;"Why, Mary loves the lamb, you know" the teacher did reply."
Input sample:

Your program should accept as its first argument a path to a filename.The input file contains several lines. Each line is one test case. Each line contains a number followed by a string, separated by a comma. eg

2,the
The first number is the n-gram length. The second string is the text printed by the user and whose prediction you have to print out.
Output sample:

For each set of input produce a single line of output which is the predictions for what the user is going to type next. eg.

lamb,0.375;teacher,0.250;children,0.125;eager,0.125;rule,0.125
Things that are good: Documentation, Readability, Simplicity. Things that are not so good: Pre-optimization

Bonus points for:
1. Making it able to dynamically add new corpus data. 
2.Allow reverse ngram lookups to find what words commonly precede words/phrases
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <map>
#include <iomanip>

using namespace std;

bool valid(char cur) {
	if (cur >= 'a' && cur <= 'z') return true;
	if (cur >= 'A' && cur <= 'Z') return true;
	if (cur >= '0' && cur <= '9') return true;
	if (cur == ' ') return true;

	return false;
}

void removeNonAlpha(string &msg) {
	string res = "";
	for (int i = 0; i < msg.size(); ++i) {
		if (valid(msg[i])) res += msg[i];
	}
	msg = res;
}

bool compare(pair<string, int> a, pair<string, int> b) {
	if (a.second > b.second) return true;
	if (a.second < b.second) return false;
	if (a.first.compare(b.first) < 0) return true;
	return false;
}

int main(int argc, char *argv[]) {
    string msg = "Mary had a little lamb its fleece was white as snow; And everywhere that Mary went, the lamb was sure to go. It followed her to school one day, which was against the rule; It made the children laugh and play, to see a lamb at school. And so the teacher turned it out, but still it lingered near, And waited patiently about till Mary did appear. \"Why does the lamb love Mary so?\" the eager children cry;\"Why, Mary loves the lamb, you know\" the teacher did reply.\"";
    
	removeNonAlpha(msg);
	transform(msg.begin(), msg.end(), msg.begin(), ::tolower);

	// vector<string> words;
	// istringstream iss(msg);
	// string temp;
	// while (iss >> temp) {
	// 	words.push_back(temp);
	// }

	ifstream ifs(argv[1]);
    string line;

	while (getline(ifs, line)) {
		line = line.substr(line.find(',')+1);
		transform(line.begin(), line.end(), line.begin(), ::tolower);
		
		map<string, int> res;
		int count = 0;
		int pos = msg.find(line);
		while (pos != string::npos) {
			int pos2 = msg.find(' ', pos+line.size()+1);
			string tmp = msg.substr(pos+line.size()+1, pos2-pos-line.size()-1);
			res[tmp]++;
			count++;
			pos = msg.find(line, pos+line.size()+1);
		}

		vector<pair<string, int> > table;
		for (map<string, int>::iterator it = res.begin(); it != res.end(); ++it) {
			table.push_back(make_pair<string, int> (it->first, it->second));
		}
		
		sort(table.begin(), table.end(), compare);
		for (int i = 0; i < table.size(); ++i) {
			cout << table[i].first << "," << fixed << setprecision(3) << (double)(table[i].second/(double)count);
			if (i != table.size()-1) cout << ";";
		}
		cout << endl;

	}

    return 0;
}
















