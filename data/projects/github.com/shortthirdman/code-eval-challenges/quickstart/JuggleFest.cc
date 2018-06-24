/*
Juggle Fest  Share on LinkedIn

Description:

Many developers here at Yodle are avid jugglers. To celebrate their prowess we are organizing a Yodle Open JuggleFest, but we need your help planning it. There will be thousands of participants split into teams. Each team will attempt to complete a juggling circuit consisting of several tricks. Each circuit emphasizes different aspects of juggling, requiring hand to eye coordination (H), endurance (E) and pizzazz (P) in various amounts for successful completion. Each juggler has these abilities in various amounts as well. How good a match they are for a circuit is determined by the dot product of the juggler's and the circuit's H, E, and P values. The higher the result, the better the match.

Each participant will be on exactly one team and there will be a distinct circuit for each team to attempt. Each participant will rank in order of preference their top X circuits. Since we would like the audiences to enjoy the performances as much as possible, when assigning jugglers to circuits we also want to consider how well their skills match up to the circuit. In fact we want to match jugglers to circuits such that no juggler could switch to a circuit that they prefer more than the one they are assigned to and be a better fit for that circuit than one of the other jugglers assigned to it.

To help us create the juggler/circuit assignments write a program in a language of your choice that takes as input a file of circuits and jugglers and outputs a file of circuits and juggler assignments. The number of jugglers assigned to a circuit should be the number of jugglers divided by the number of circuits. Assume that the number of circuits and jugglers will be such that each circuit will have the same number of jugglers with no remainder.
Input sample:

One line per circuit or juggler. All circuits will come before any jugglers. Circuit lines start with a C and juggler lines start with a J. Names of circuits and jugglers will never have spaces. A skill and the rating for that skill are separated by a colon. Circuit lines have the circuit names followed by skills. Juggler lines have the juggler names followed by skills, followed by circuits in order of preference, separated by commas. Example:

C C0 H:7 E:7 P:10
C C1 H:2 E:1 P:1
C C2 H:7 E:6 P:4

J J0 H:3 E:9 P:2 C2,C0,C1
J J1 H:4 E:3 P:7 C0,C2,C1
J J2 H:4 E:0 P:10 C0,C2,C1
J J3 H:10 E:3 P:8C2,C0,C1
J J4 H:6 E:10 P:1 C0,C2,C1
J J5 H:6 E:7 P:7 C0,C2,C1
J J6 H:8 E:6 P:9 C2,C1,C0
J J7 H:7 E:1 P:5 C2,C1,C0
J J8 H:8 E:2 P:3 C1,C0,C2
J J9 H:10 E:2 P:1 C1,C2,C0
J J10 H:6 E:4 P:5 C0,C2,C1
J J11 H:8 E:4 P:7 C0,C1,C2
Output sample:

You should have one line per circuit assignment. Each line should contain the circuit name followed by the juggler name, followed by that juggler's circuits in order of preference and the match score for that circuit. The line should include all jugglers matched to the circuit. The example below is a valid assignment for the input file above (has 3 lines).

C2 J6 C2:128 C1:31 C0:188, J3 C2:120 C0:171 C1:31, J10 C0:120 C2:86 C1:21, J0 C2:83 C0:104 C1:17

C1 J9 C1:23 C2:86 C0:94, J8 C1:21 C0:100 C2:80, J7 C2:75 C1:20 C0:106, J1 C0:119 C2:74 C1:18

C0 J5 C0:161 C2:112 C1:26, J11 C0:154 C1:27 C2:108, J2 C0:128 C2:68 C1:18, J4 C0:122 C2:106 C1:23
Run your program on the input which contains 2000 circuits and 12000 jugglers. The correct output is the sum of the names of the jugglers (taking off the leading letter J) that are assigned to circuit C1970. So for example if the jugglers assigned to circuit C1970 were J1,J2,J3,J4,J5 and J6 the correct answer would be
21
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <set>

using namespace std;

typedef struct Juggler {
	string name;
	int H, E, P;
	vector<string> prefer;
	string team;
	Juggler(string n, int h, int e, int p) : name(n), H(h), E(e), P(p) {}
} Juggler;

typedef struct Circle {
	string name;
	int H, E, P;
	vector<Juggler *> member;
	Circle(string n, int h, int e, int p) : name(n), H(h), E(e), P(p) {}
} Circle;

void parse(vector<Circle> &c, vector<Juggler> &j, string buf) {
	char T = buf[0];
	buf = buf.substr(2);  // trim type

	string name = buf.substr(0, buf.find(' '));
	buf = buf.substr(buf.find(' ')+1);	// trim name

	string H = buf.substr(0, buf.find(' '));
	buf = buf.substr(buf.find(' ')+1);	// trim H
	int Hval = atoi(H.substr(2).c_str());

	string E = buf.substr(0, buf.find(' '));
	buf = buf.substr(buf.find(' ')+1);	// trim H
	int Eval = atoi(E.substr(2).c_str());

	string P = buf.substr(0, buf.find(' '));
	buf = buf.substr(buf.find(' ')+1);	// trim H
	int Pval = atoi(P.substr(2).c_str());

	if (T == 'C') {
		struct Circle curCircle(name, Hval, Eval, Pval);
		c.push_back(curCircle);
	} else {
		struct Juggler curJuggler(name, Hval, Eval, Pval);
		int pos = buf.find(',');
		while (pos != string::npos) {
			curJuggler.prefer.push_back(buf.substr(0, pos));
			buf = buf.substr(pos+1);
			pos = buf.find(',');
		}
		curJuggler.prefer.push_back(buf);
		j.push_back(curJuggler);
	}
}

int findCircle(vector<Circle> &c, string s) {
	for (int i = 0; i < c.size(); ++i) {
		if (c[i].name == s) return i;
	}
	return -1;
}

void assign(Circle &c, Juggler &j) {
	j.team = c.name;
	c.member.push_back(&j);
}

void initAssign(vector<Circle> &circles, vector<Juggler> &jugglers) {
	int num = jugglers.size() / circles.size();

	for (int i = 0; i < jugglers.size(); ++i) {
		int j = 0;
		while (j < jugglers[i].prefer.size()) {
			int cur = findCircle(circles, jugglers[i].prefer[j]);
			if (cur != -1 && circles[cur].member.size() < num) {
				assign(circles[cur], jugglers[i]);
				break;
			}
			j++;
		}
		if (j == jugglers[i].prefer.size()) {
			for (int k = 0; k < circles.size(); ++k) {
				if (circles[k].member.size() < num) {
					assign(circles[k], jugglers[i]);
					break;
				}
			}
		}
	}
}

void print(vector<Circle> &circles, vector<Juggler> &jugglers) {
	cout << "Circles:" << endl;
	for (vector<Circle>::iterator it = circles.begin(); it != circles.end(); ++it) {
		cout << it->name << " " << it->H << " " << it->E << " " << it->P;
		cout << " member: ";
		for (int j = 0; j < it->member.size(); ++j) {
			cout << it->member[j]->name << " ";
		}
		cout << endl;
	}
	cout << endl;

	cout << "Jugglers:" << endl;
	for (int i = 0; i < jugglers.size(); ++i) {
		cout << jugglers[i].name << " " << jugglers[i].H << " " << jugglers[i].E << " " << jugglers[i].P;
		cout << " prefer: ";
		for (int j = 0; j < jugglers[i].prefer.size(); ++j) {
			cout << jugglers[i].prefer[j] << " ";
		}
		cout << ", on team: " << jugglers[i].team << endl;
	}
	cout << endl;   
}

int matchVal(Circle c, Juggler j) {
	return (c.H * j.H) + (c.E * j.E) + (c.P * j.P);
}

void optimize(vector<Circle> &circles, vector<Juggler> &jugglers) {
	while (1) {
		bool change = false;

		for (int i = 0; i < jugglers.size(); ++i) {
			bool found = false;

			for (int j = 0; j < jugglers[i].prefer.size(); ++j) {
				if (jugglers[i].team == jugglers[i].prefer[j]) break;
				
				int dst = findCircle(circles, jugglers[i].prefer[j]);
				for (int k = 0; k < circles[dst].member.size(); ++k) {

					Juggler * tmp = circles[dst].member[k];

					if (matchVal(circles[dst], jugglers[i]) > matchVal(circles[dst], *tmp)) {
						change = true;  found = true;

						int src = findCircle(circles, jugglers[i].team);

						circles[src].member.erase(find(circles[src].member.begin(), circles[src].member.end(), &jugglers[i]));
						circles[dst].member.erase(find(circles[dst].member.begin(), circles[dst].member.end(), tmp));

						assign(circles[src], *tmp);
						assign(circles[dst], jugglers[i]);

						break;
					}
				}

				if (found) break;
			}
		}

		if (!change) return;
	}
}

int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string buf;
    vector<Circle> circles;
    vector<Juggler> jugglers;
    
    while (getline(ifs, buf)) {
    	if (buf.size() == 0) continue;
    	parse(circles, jugglers, buf);
    }

    initAssign(circles, jugglers);

    // print(circles, jugglers);

    optimize(circles, jugglers);

    // print(circles, jugglers);

    int res = 0;
    for (int i = 0; i < circles.size(); ++i) {
    	if (circles[i].name == "C1970") {
    		for (int j = 0; j < circles[i].member.size(); ++j) {
    			Juggler * cur = circles[i].member[j];
    			res += atoi(cur->name.substr(1).c_str());
    		}
    	}
    }

    cout << res << endl;

    return 0;
}