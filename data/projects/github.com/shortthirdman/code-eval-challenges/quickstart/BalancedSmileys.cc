/*
Balanced Smileys  Share on LinkedIn

Description:

Credits: This problem appeared in the Facebook Hacker Cup 2013 Hackathon.

Your friend John uses a lot of emoticons when you talk to him on Messenger. In addition to being a person who likes to express himself through emoticons, he hates unbalanced parenthesis so much that it makes him go :(.

Sometimes he puts emoticons within parentheses, and you find it hard to tell if a parenthesis really is a parenthesis or part of an emoticon. A message has balanced parentheses if it consists of one of the following:

- An empty string ""
- One or more of the following characters: 'a' to 'z', ' ' (a space) or ':' (a colon)
- An open parenthesis '(', followed by a message with balanced parentheses, followed by a close parenthesis ')'.
- A message with balanced parentheses followed by another message with balanced parentheses.
- A smiley face ":)" or a frowny face ":("

Write a program that determines if there is a way to interpret his message while leaving the parentheses balanced.

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file contains a message that you got from John. e.g.

:((
i am sick today (:()
(:)
hacker cup: started :):)
)(
Output sample:

Print out the string "YES"/"NO"(all quotes for clarity only) stating whether or not it is possible that the message had balanced parentheses. e.g.

NO
YES
YES
YES
NO
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <stack>

using namespace std;

bool valid(string s, int cur) {
	if (cur < 0) return false;
	if (s.size() == 0) return cur == 0;
	if (s.size() == 1) {
		if (s[0] == ')') return cur == 1;
		else if (s[0] == '(') return false;
		else return cur == 0;
	}

	if (s[0] == '(') { 
		return valid(s.substr(1), cur+1);
	} else if (s[0] == ')') {
		return valid(s.substr(1), cur-1);
	} else if (s[0] == ':' && (s[1] == '(' || s[1] == ')')) {
		return valid(s.substr(1), cur) || valid(s.substr(2), cur);
	} else {
		return valid(s.substr(1), cur);
	}
}

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    
    while (getline(ifs, line)) {
        if (valid(line, 0)) {
        	cout << "YES" << endl;
        } else {
        	cout << "NO" << endl;
        }
    }

    return 0;
}