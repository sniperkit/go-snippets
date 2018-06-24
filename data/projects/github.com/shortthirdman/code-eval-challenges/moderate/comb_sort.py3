import fileinput

SHRINK = 1.25

for line in fileinput.input():
    s = [int(i) for i in line.split()]
    n = 0
    gap = len(s)
    while not all(s[i] <= s[i+1] for i in range(len(s)-1)):
        n += 1
        if gap > 1:
            gap = int(gap / SHRINK)
        for i in range(len(s) - gap):
            if s[i] > s[i + gap]:
                s[i], s[i + gap] = s[i + gap], s[i]
    print(n)
