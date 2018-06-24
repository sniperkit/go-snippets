"""Count the number of prime numbers between N and M (both inclusive)."""
import fileinput

def isp(num, pri):
    """Test for primality."""
    for j in pri:
        if num % j == 0:
            return False
    return True

PRIMES = [2, 3]

for line in fileinput.input():
    m, n = [int(i) for i in line.split(',')]
    if m > n:
        m, n = n, m
    i = PRIMES[-1]
    while i <= n:
        if isp(i, PRIMES):
            PRIMES += [i]
        i += 2
    i, f = 0, False
    while PRIMES[i] < m:
        i += 1
        if i == len(PRIMES) or PRIMES[i] > n:
            print(0)
            f = True
            break
    if f:
        continue
    c = 1
    i += 1
    while i < len(PRIMES) and PRIMES[i] <= n:
        i += 1
        c += 1
    print(c)
