import fileinput

for line in fileinput.input():
    n, a, b = [int(i) for i in line.split(',')]
    print("true" if (n & (1 << (a - 1)) == 0) == (n & (1 << (b - 1)) == 0) else "false")
