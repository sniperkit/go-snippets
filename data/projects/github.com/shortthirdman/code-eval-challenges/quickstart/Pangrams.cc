/*
Pangrams  Share on LinkedIn

Description:

The sentence 'A quick brown fox jumps over the lazy dog' contains every single letter in the alphabet. Such sentences are called pangrams. You are to write a program, which takes a sentence, and returns all the letters it is missing (which prevent it from being a pangram). You should ignore the case of the letters in sentence, and your return should be all lower case letters, in alphabetical order. You should also ignore all non US-ASCII characters.In case the input sentence is already a pangram, print out the string NULL

Input sample:

Your program should accept as its first argument a filename. This file will contain several text strings, one per line. Ignore all empty lines. eg.

A quick brown fox jumps over the lazy dog
A slow yellow fox crawls under the proactive dog
Output sample:

Print out all the letters each string is missing in lowercase, alphabetical order .e.g.

NULL
bjkmqz
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <algorithm>
#include <vector>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    vector<int> table;

    while (getline(ifs, line)) {
    	table.clear();
    	table.resize(26);
    	transform(line.begin(), line.end(), line.begin(), ::tolower);
        for (int i = 0; i < line.size(); ++i) {
        	if (line[i] >= 'a' && line[i] <= 'z') table[line[i]-'a']++;
        }

        int found = 0;
        for (int i = 0; i < table.size(); ++i) {
        	if (table[i] == 0) {
        		found = 1;
        		cout << (char)(i+'a');
        	}
        }
        if (found == 0) { cout << "NULL"; }
        cout << endl;
    }

    return 0;
}