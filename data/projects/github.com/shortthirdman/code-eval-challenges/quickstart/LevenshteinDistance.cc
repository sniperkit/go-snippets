/*
Levenshtein Distance  Share on LinkedIn

Description:

Two words are friends if they have a Levenshtein distance of 1 (For details see http://en.wikipedia.org/wiki/Levenshtein_distance). That is, you can add, remove, or substitute exactly one letter in word X to create word Y. A word’s social network consists of all of its friends, plus all of their friends, and all of their friends’ friends, and so on. Write a program to tell us how big the social network for the word 'hello' is, using this word list https://raw.github.com/codeeval/Levenshtein-Distance-Challenge/master/input_levenshtein_distance.txt

Input sample:

Your program should accept as its first argument a path to a filename.The input file contains the word list. This list is also available at https://raw.github.com/codeeval/Levenshtein-Distance-Challenge/master/input_levenshtein_distance.txt.

Output sample:

Print out how big the social network for the word 'hello' is. e.g. The social network for the word 'abcde' is 4846.
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <set>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    set<string> dict;
    vector<string> res;
    res.push_back("hello");

    while (getline(ifs, line)) {
        if (line == "hello") continue;
        dict.insert(line);
    }

	for (int i = 0; i < res.size(); ++i) {
		for (int j = 0; j < res[i].size(); ++j) {
			string cur;

			// substitution
			cur = res[i];
			for (int k = 0; k < 26; ++k) {
				cur[j] = 'a'+k;
				set<string>::iterator it = dict.find(cur);
				if (it != dict.end()) {
					dict.erase(it);
					res.push_back(cur);
				}
			}

			// add
			for (int k = 0; k <= res[i].size(); ++k) {
				cur = res[i];
				cur.insert(k, "a");
				for (int x = 0; x < 26; ++x) {
					cur[k] = 'a'+x;
					set<string>::iterator it = dict.find(cur);
					if (it != dict.end()) {
						dict.erase(it);
						res.push_back(cur);
					}
				}
			}

			// delete
			for (int k = 0; k < res[i].size(); ++k) {
				cur = res[i];
				cur.erase(k, 1);
				set<string>::iterator it = dict.find(cur);
				if (it != dict.end()) {
					dict.erase(it);
					res.push_back(cur);
				}
			}
		}
	}

	cout << res.size()-1 << endl;

    return 0;
}