/*
Poker hands  Share on LinkedIn

Description:

In the card game poker, a hand consists of five cards and are ranked, from lowest to highest, in the following way: 

High Card: Highest value card. 
One Pair: Two cards of the same value. 
Two Pairs: Two different pairs. 
Three of a Kind: Three cards of the same value. 
Straight: All cards are consecutive values. 
Flush: All cards of the same suit. 
Full House: Three of a kind and a pair. 
Four of a Kind: Four cards of the same value. 
Straight Flush: All cards are consecutive values of same suit. 
Royal Flush: Ten, Jack, Queen, King, Ace, in same suit. 
The cards are valued in the order: 
2, 3, 4, 5, 6, 7, 8, 9, Ten, Jack, Queen, King, Ace. 

If two players have the same ranked hands then the rank made up of the highest value wins; for example, a pair of eights beats a pair of fives. But if two ranks tie, for example, both players have a pair of queens, then highest cards in each hand are compared; if the highest cards tie then the next highest cards are compared, and so on. 
Input sample:

Your program should accept as its first argument a path to a filename. Each line in this file contains 2 hands (left and right). Cadrs and hands are separated by space. E.g.

6D 7H AH 7S QC 6H 2D TD JD AS
JH 5D 7H TC JS JD JC TS 5S 7S
2H 8C AD TH 6H QD KD 9H 6S 6C
JS JH 4H 2C 9H QH KC 9D 4D 3S
TC 7H KH 4H JC 7D 9S 3H QS 7S
Output sample:

Print out the name of the winning hand or "none" in case the hands are equal. E.g.

left
none
right
left
right
*/

#include <iostream>
#include <fstream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <map>

using namespace std;

typedef pair<int, char> Card;
typedef bool (*FunctionPointer)(vector<Card> &, int &);

Card parse(string c) {
    static string values = "23456789TJQKA";
    return Card(values.find(c[0])+2, c[1]);
}

bool byValue(Card a, Card b) {
    return a.first > b.first;
}

int totalVal(vector<Card> &a) {
    int val = 0;
    for (int i = 0; i < a.size(); ++i)
        val += a[i].first;
    return val;
}


// High Card: Highest value card. 
bool highCard(vector<Card> &a, int & high) {
    high = a[0].first;

    return true;
}

// One Pair: Two cards of the same value. 
bool onePair(vector<Card> &a, int & high) {
    for (int i = 0; i < a.size()-1; ++i) {
        if (a[i].first == a[i+1].first) {
            high = a[i].first;
            return true;
        }
    }
    return false;
}

// Two Pairs: Two different pairs. 
// Note: here high has to be the sum of two pairs value.
bool twoPairs(vector<Card> &a, int & high) {
    int pair = 0;
    high = 0;
    for (int i = 0; i < a.size()-1; ++i) {
        if (a[i].first == a[i+1].first) {
            high += a[i].first;
            ++i; ++pair;
        }
    }
    
    return (pair == 2);
}

// Three of a Kind: Three cards of the same value. 
bool threeOfAKind(vector<Card> &a, int & high) {
    for (int i = 0; i < a.size()-2; ++i) {
        if (a[i].first == a[i+1].first && a[i].first == a[i+2].first) {
            high = a[i].first;
            return true;
        }
    }
    return false;
}

// Straight: All cards are consecutive values. 
bool straight(vector<Card> &a, int & high) {
    for (int i = 1; i < a.size(); ++i) {
        if (a[i].first != a[0].first-i) return false;
    }
    high = a[0].first;
    return true;
}

// Flush: All cards of the same suit. 
bool flush(vector<Card> &a, int & high) {
    for (int i = 1; i < a.size(); ++i) {
        if (a[i].second != a[0].second) return false;
    }
    high = a[0].first;
    return true;
}

// Full House: Three of a kind and a pair. 
// card is sorted by value
bool fullHouse(vector<Card> &a, int & high) {
    bool house23 = (a[0].first == a[1].first && a[2].first == a[3].first && a[3].first == a[4].first);
    bool house32 = (a[0].first == a[1].first && a[1].first == a[2].first && a[3].first == a[4].first);

    if (house23 || house32) {
        if (house23) high = a[0].first;
        else         high = a[4].first;
        return true;
    }

    return false;
}

// Four of a Kind: Four cards of the same value. 
bool fourOfAKind(vector<Card> &a, int & high) {
    map<int, int> exist;
    for (int i = 1; i < a.size(); ++i)
        exist[a[i].first]++;

    if (exist[a[0].first] == 4) {
        high = a[0].first;
        return true;
    }
    if (exist[a[4].first] == 4) {
        high = a[4].first;
        return true;
    }

    return false;
}

// Straight Flush: All cards are consecutive values of same suit. 
bool straightFlush(vector<Card> &a, int & high) {
    for (int i = 1; i < a.size(); ++i) {
        if (a[i].second != a[0].second) return false;
        if (a[i].first != a[0].first-i) return false;
    }
    high = a[0].first;

    return true;
}

// Royal Flush: Ten, Jack, Queen, King, Ace, in same suit. 
bool royalFlush(vector<Card> &a, int & high) {
    if (a[0].first != 14) return false;

    return straightFlush(a, high);
}

string compare(vector<Card> &left, vector<Card> &right) {
    static FunctionPointer fps[] = {highCard, onePair, twoPairs, threeOfAKind, straight, flush, fullHouse, fourOfAKind, straightFlush, royalFlush};
    int nRanks = 10;

    sort(left.begin(), left.end(), byValue);
    int leftVal = 0, leftRank = 0;
    for (leftRank = nRanks-1; leftRank >= 0; --leftRank)
        if (fps[leftRank](left, leftVal))   break;

    sort(right.begin(), right.end(), byValue);
    int rightVal = 0, rightRank = 0;
    for (rightRank = nRanks-1; rightRank >= 0; --rightRank)
        if (fps[rightRank](right, rightVal))   break;
    
    // cout << "leftRank: " << leftRank << ", leftVal: " << leftVal << endl;
    // cout << "rightRank: " << rightRank << ", rightVal: " << rightVal << endl;

    if (leftRank != rightRank) {
        if (leftRank < rightRank) return "right";
        else                      return "left";
    } else if (leftVal != rightVal) {
        if (leftVal < rightVal) return "right";
        else                    return "left";
    } else {
        leftVal = totalVal(left);
        rightVal = totalVal(right);

        if (leftVal < rightVal) return "right";
        else if (leftVal == rightVal) return "none";
        else return "left";
    }
}


int main(int argc, char *argv[]) {
    ifstream ifs(argv[1]);
    string line;
    
    while (getline(ifs, line)) {
        istringstream iss(line);
        vector<Card> left;
        vector<Card> right;

        string buf;
        while (iss >> buf) {
            if (left.size() < 5) left.push_back(parse(buf));
            else                right.push_back(parse(buf));
        }

        // cout << "left : ";
        // for (int i = 0; i < left.size(); ++i)
        //     cout << left[i].first << left[i].second << " ";
        // cout << endl;

        // cout << "right: ";
        // for (int i = 0; i < right.size(); ++i)
        //     cout << right[i].first << right[i].second << " ";
        // cout << endl;

        cout << compare(left, right) << endl;
    }

    ifs.close();
    return 0;
}