import fileinput

for line in fileinput.input():
    s = line.split(" | ")
    t = [int(i) for i in s[0].split()]
    n = min(int(s[1]), len(t)-1)
    for i in range(n):
        for j in range(1, len(t)):
            if t[j] < t[j-1]:
                t[j], t[j-1] = t[j-1], t[j]
    print(*t)
