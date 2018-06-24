import fileinput

for line in fileinput.input():
    s = line.split(';')
    t = [int(i) for i in s[1].split()]
    n, m = int(s[0]), 0
    c = sum(t[:n])
    if c > m:
        m = c
    for i in range(n, len(t)):
        c = c - t[i-n] + t[i]
        if c > m:
            m = c
    print(m)
