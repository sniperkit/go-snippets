/*
Sum of Primes  Share on LinkedIn

Description:

Write a program to determine the sum of the first 1000 prime numbers.
Input sample:

None

Output sample:

Your program should print the sum on stdout.i.e.

3682913
*/

#include <iostream>

using namespace std;

bool prime(long long target) {
	if (target == 2) return true;
	if (target % 2 == 0) return false;
    for (long long i = 3; i < target/2; i += 2) {
        if (target % i == 0) return false;
    }
    return true;
}

int main(int argc, char *argv[]) {
	long long sum = 0;

	int count = 0;
	long long cur = 2;
	while (count < 1000) {
		if (prime(cur)) {
			sum += cur;
			++count;
		}
		++cur;
	}

	cout << sum << endl;

    return 0;
}