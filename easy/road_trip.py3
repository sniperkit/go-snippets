import fileinput

for line in fileinput.input():
    s = [int(i) for i in ''.join(c for c in line if c in "0123456789 ").split()]
    s.sort()
    for i in range(len(s)-1, 0, -1):
        s[i] -= s[i-1]
    print(*s, sep=',')
