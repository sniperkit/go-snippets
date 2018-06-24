import fileinput

for line in fileinput.input():
    s = line.split(" | ")
    t = [int(i) for i in s[0].split()]
    n = min(int(s[1]), len(t)//2)
    for i in range(n):
        for j in range(i + 1, len(t) - i):
            if t[j-1] > t[j]:
                t[j-1], t[j] = t[j], t[j-1]
        for j in range(len(t) - i - 1, i, -1):
            if t[j-1] > t[j]:
                t[j-1], t[j] = t[j], t[j-1]
    print(' '.join([str(i) for i in t]))
