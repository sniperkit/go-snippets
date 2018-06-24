/*
Prime Palindrome  Share on LinkedIn

Description:

Write a program to determine the biggest prime palindrome under 1000.
Input sample:

None

Output sample:

Your program should print the largest palindrome on stdout. i.e.

929

Submit your solution in a file (some file name).(py| c| cpp| rb| pl| php| tcl| clj| js| scala| cs| m) | prime_palindrome.java|prime_palindrome.scala or use the online editor.
*/

#include <iostream>

using namespace std;

bool prime(int target) {
    for (int i = 2; i < target/2; ++i) {
        if (target % i == 0) return false;
    }
    return true;
}

bool palin(int target) {
    if (target >= 100) {
        if (target/100 == target%10) return true;
    } else if (target >= 10) {
        if (target/10 == target%10) return true;
    } else {
        return true;
    }
    return false;
}

int main(int argc, char *argv[]) {
    for (int i = 1000; i >= 0; --i) {
        if (prime(i) && palin(i)) {
            cout << i << endl;
            return 0;
        }
    }

    return 0;
}