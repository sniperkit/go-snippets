/*
URI Comparison  Share on LinkedIn

Description:

Determine if two URIs match. For the purpose of this challenge, you should use a case-sensitive octet-by-octet comparison of the entire URIs, with these exceptions:

1. A port that is empty or not given is equivalent to the default port of 80
2. Comparisons of host names MUST be case-insensitive
3. Comparisons of scheme names MUST be case-insensitive
4. Characters are equivalent to their % HEX HEX encodings. (Other than typical reserved characters in urls like /,?,@,:,&,= etc)

Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file contains two urls delimited by a semicolon. e.g.

http://abc.com:80/~smith/home.html;http://ABC.com/%7Esmith/home.html
Output sample:

Print out True/False if the URIs match. e.g.

True

*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <string>

using namespace std;

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    
    while (getline(ifs, line)) {
    	transform(line.begin(), line.end(), line.begin(), ::tolower);

    	int pos = line.find(';');
    	string left = line.substr(0, pos);
    	string right = line.substr(pos+1);

        string w_port = ".com:80";
        string wo_port = ".com";
        pos = left.find(w_port);
        if (pos != string::npos) left.replace(pos, w_port.size(), wo_port);
        pos = right.find(w_port);
        if (pos != string::npos) right.replace(pos, w_port.size(), wo_port);

        w_port = ".com:";
        pos = left.find(w_port);
        if (pos != string::npos) left.replace(pos, w_port.size(), wo_port);
        pos = right.find(w_port);
        if (pos != string::npos) right.replace(pos, w_port.size(), wo_port);

        pos = left.find('%');
        while (pos != string::npos) {
            string h = left.substr(pos+1, 2);
            const char *c = h.c_str();            
            unsigned int x;
            sscanf(c, "%2X", &x);
            string tmp;
            tmp += x;
            left.replace(pos, 3, tmp);
            pos = left.find('%');
        }

        pos = right.find('%');
        while (pos != string::npos) {
            string h = right.substr(pos+1, 2);
            const char *c = h.c_str();            
            unsigned int x;
            sscanf(c, "%2X", &x);
            string tmp;
            tmp += x;
            right.replace(pos, 3, tmp);
            pos = right.find('%');
        }

        if (left.compare(right) == 0) {
            cout << "True" << endl;
        } else {
            cout << "False" << endl;
        }

    }

    return 0;
}