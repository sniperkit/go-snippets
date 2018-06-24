/*
Decryption

Description:

For this challenge you are given an encrypted message and a key. You have to determine the encryption and encoding technique and then print out the corresponding plaintext message. You can assume that the plaintext corresponding to this message, and all messages you must handle, will be comprised of only the characters A-Z and spaces; no digits or punctuation.

Input sample:

There is no input for this program. The encrypted message and key is:

message: "012222 1114142503 0313012513 03141418192102 0113 2419182119021713 06131715070119",
keyed_alphabet: "BHISOECRTMGWYVALUZDNFJKPQX"
Output sample:

Print out the plaintext message. (in CAPS)
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <map>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string KEY = "BHISOECRTMGWYVALUZDNFJKPQX";
    string MSG = "012222 1114142503 0313012513 03141418192102 0113 2419182119021713 06131715070119";

    map<char, int> map;
    for (int i = 0; i < KEY.size(); ++i) {
    	map[KEY[i]] = i;
    }

    int i = 0;
    while (i < MSG.size()) {
    	if (MSG[i] == ' ') {
    		cout << " ";
    		i++; continue;
    	}
    	int val = 0;
    	val = MSG[i++]-'0';
    	val *= 10;
    	val += MSG[i++]-'0';

    	cout << (char)('A'+map[(char)('A'+val)]);
    }
    cout << endl;

    return 0;
}