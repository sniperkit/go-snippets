import fileinput

for line in fileinput.input():
    n, a = [int(i) for i in line.split()]
    r = 0
    for i in range(1 << n, a + 1):
        m, b = n, i
        while b > 0 and m >= 0:
            if b % 2 == 0:
                m -= 1
            b >>= 1
        if m == 0:
            r += 1
    print(r)
